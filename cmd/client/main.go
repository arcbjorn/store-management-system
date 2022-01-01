package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/arcbjorn/store-management-system/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient laptop.LaptopServiceClient, lp *laptop.Laptop) {
	req := &laptop.CreateLaptopRequest{
		Laptop: lp,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}

	log.Printf("created laptop with id: %s", res.Id)
}

func searchLaptop(laptopClient laptop.LaptopServiceClient, filter *laptop.Filter) {
	log.Print("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &laptop.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}

		lp := res.GetLaptop()
		log.Print("- found: ", lp.GetId())
		log.Print("  - brand: ", lp.GetBrand())
		log.Print("  - name: ", lp.GetName())
		log.Print("  - cpu cores: ", lp.GetCpu().GetCoreNumber())
		log.Print("  - cpu min ghz: ", lp.GetCpu().GetMinGhz())
		log.Print("  - ram: ", lp.GetRam().GetValue(), lp.GetRam().GetUnit())
		log.Print("  - price: ", lp.GetPriceUsd())
	}
}

func uploadImage(laptopClient laptop.LaptopServiceClient, laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &laptop.UploadImageRequest{
		Data: &laptop.UploadImageRequest_Info{
			Info: &laptop.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)

	if err != nil {
		log.Fatal("cannot send image info: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &laptop.UploadImageRequest{
			Data: &laptop.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

func testCreateLaptop(laptopClient laptop.LaptopServiceClient) {
	createLaptop(laptopClient, sample.NewLaptop())
}

func testSearchLaptop(laptopClient laptop.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
	}

	filter := &laptop.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &laptop.Memory{Value: 8, Unit: laptop.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
}

func testUploadImage(laptopClient laptop.LaptopServiceClient) {
	lp := sample.NewLaptop()
	createLaptop(laptopClient, lp)
	uploadImage(laptopClient, lp.GetId(), "tmp/laptop.jpg")
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial sever %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := laptop.NewLaptopServiceClient(conn)
	testUploadImage(laptopClient)
}
