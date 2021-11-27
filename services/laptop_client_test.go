package services_test

import (
	"context"
	"net"
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

	laptopServer, serverAddress := startTestLaptopServer(t)
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

	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, res)

	requireSameLaptop(t, newLaptop, other)
}

func startTestLaptopServer(t *testing.T) (*services.LaptopServer, string) {
	laptopServer := services.NewLaptopServer(services.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	laptop.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
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
