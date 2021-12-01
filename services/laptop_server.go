package services

import (
	"context"
	"errors"
	"log"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server that provides services for laptop functionality
type LaptopServer struct {
	laptop.UnimplementedLaptopServiceServer
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *laptop.CreateLaptopRequest,
) (*laptop.CreateLaptopResponse, error) {
	laptopDto := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptopDto.Id)

	if len(laptopDto.Id) > 0 {
		_, err := uuid.Parse(laptopDto.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptopDto.Id = id.String()
	}

	if ctx.Err() == context.Canceled {
		log.Print("request is cancelled")
		return nil, status.Error(codes.Canceled, "request is cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	// save new Laptop to store
	err := server.Store.Save(laptopDto)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save latop to the store: %v", err)
	}

	log.Printf("save laptop with id: %s", laptopDto.Id)

	res := &laptop.CreateLaptopResponse{
		Id: laptopDto.Id,
	}
	return res, nil
}

func (server *LaptopServer) SearchLaptop(
	req *laptop.SearchLaptopRequest,
	stream laptop.LaptopService_SearchLaptopServer,
) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := server.Store.Search(
		filter,
		func(lp *laptop.Laptop) error {
			res := &laptop.SearchLaptopResponse{Laptop: lp}

			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", lp.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}
