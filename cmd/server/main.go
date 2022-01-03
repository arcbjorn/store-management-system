package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/arcbjorn/store-management-system/services"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	laptopStore := services.NewInMemoryLaptopStore()
	imageStore := services.NewDiskImageStore("img")
	ratingStore := services.NewInMemoryRatingStore()

	laptopServer := services.NewLaptopServer(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer()
	laptop.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
