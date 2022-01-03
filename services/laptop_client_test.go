package services_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/arcbjorn/store-management-system/sample"
	"github.com/arcbjorn/store-management-system/serializer"
	"github.com/arcbjorn/store-management-system/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := services.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	newLaptop := sample.NewLaptop()
	expectedID := newLaptop.Id
	req := &laptop.CreateLaptopRequest{
		Laptop: newLaptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, res)

	requireSameLaptop(t, newLaptop, other)
}

func startTestLaptopServer(
	t *testing.T,
	laptopStore services.LaptopStore,
	imageStore services.ImageStore,
	ratingStore services.RatingStore,
) string {
	laptopServer := services.NewLaptopServer(laptopStore, imageStore, ratingStore)

	grpcServer := grpc.NewServer()
	laptop.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) laptop.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return laptop.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, lt1 *laptop.Laptop, lt2 *laptop.Laptop) {
	json1, err := serializer.ProtobufToJsonBytes(lt1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJsonBytes(lt2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &laptop.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &laptop.Memory{Value: 8, Unit: laptop.Memory_GIGABYTE},
	}

	laptopStore := services.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		lp := sample.NewLaptop()

		switch i {
		case 0:
			lp.PriceUsd = 2500
		case 1:
			lp.Cpu.CoreNumber = 2
		case 2:
			lp.Cpu.MinGhz = 2.0
		case 3:
			lp.Ram = &laptop.Memory{Value: 4096, Unit: laptop.Memory_MEGABYTE}
		case 4:
			lp.PriceUsd = 1999
			lp.Cpu.CoreNumber = 4
			lp.Cpu.MinGhz = 2.5
			lp.Cpu.MaxGhz = lp.Cpu.MinGhz + 2.0
			lp.Ram = &laptop.Memory{Value: 16, Unit: laptop.Memory_GIGABYTE}
			expectedIDs[lp.Id] = true
		case 5:
			lp.PriceUsd = 2000
			lp.Cpu.CoreNumber = 6
			lp.Cpu.MinGhz = 2.8
			lp.Cpu.MaxGhz = lp.Cpu.MinGhz + 2.0
			lp.Ram = &laptop.Memory{Value: 64, Unit: laptop.Memory_GIGABYTE}
			expectedIDs[lp.Id] = true
		}

		err := laptopStore.Save(lp)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &laptop.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())

		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../tmp"

	laptopStore := services.NewInMemoryLaptopStore()
	imageStore := services.NewDiskImageStore(testImageFolder)

	lp := sample.NewLaptop()
	err := laptopStore.Save(lp)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &laptop.UploadImageRequest{
		Data: &laptop.UploadImageRequest_Info{
			Info: &laptop.ImageInfo{
				LaptopId:  lp.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &laptop.UploadImageRequest{
			Data: &laptop.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))
}

func TestClientRateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := services.NewInMemoryLaptopStore()
	ratingStore := services.NewInMemoryRatingStore()

	lp := sample.NewLaptop()
	err := laptopStore.Save(lp)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, nil, ratingStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	stream, err := laptopClient.RateLaptop(context.Background())
	require.NoError(t, err)

	scores := []float64{8, 7.5, 10}
	avarages := []float64{8, 7.75, 8.5}

	n := len(scores)
	for i := 0; i < n; i++ {
		req := &laptop.RateLaptopRequest{
			LaptopId: lp.GetId(),
			Score:    scores[i],
		}

		err := stream.Send(req)
		require.NoError(t, err)
	}

	err = stream.CloseSend()
	require.NoError(t, err)

	for idx := 0; ; idx++ {
		res, err := stream.Recv()
		if err == io.EOF {
			require.Equal(t, n, idx)
			return
		}

		require.NoError(t, err)
		require.Equal(t, lp.GetId(), res.GetLaptopId())
		require.Equal(t, uint32(idx+1), res.GetRatedCount())
		require.Equal(t, avarages[idx], res.GetAverageScore())
	}
}
