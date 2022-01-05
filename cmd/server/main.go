package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/arcbjorn/store-management-system/pb/auth"
	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/arcbjorn/store-management-system/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func createUser(userStore services.UserStore, username, password, role string) error {
	user, err := services.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func seedUsers(userStore services.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/store.management.system.LaptopService"

	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	userStore := services.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	jwtManager := services.NewJWTManager(secretKey, tokenDuration)
	authServer := services.NewAuthServer(userStore, jwtManager)

	laptopStore := services.NewInMemoryLaptopStore()
	imageStore := services.NewDiskImageStore("img")
	ratingStore := services.NewInMemoryRatingStore()

	laptopServer := services.NewLaptopServer(laptopStore, imageStore, ratingStore)

	authInterceptor := services.NewAuthInterceptor(jwtManager, accessibleRoles())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	auth.RegisterAuthServiceServer(grpcServer, authServer)
	laptop.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

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
