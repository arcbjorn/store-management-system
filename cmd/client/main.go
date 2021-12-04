package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/arcbjorn/store-management-system/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient laptop.LaptopServiceClient) {
	newLaptop := sample.NewLaptop()
	req := &laptop.CreateLaptopRequest{
		Laptop: newLaptop,
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

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial sever %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := laptop.NewLaptopServiceClient(conn)

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &laptop.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &laptop.Memory{Value: 8, Unit: laptop.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
}
