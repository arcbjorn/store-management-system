package services

import (
	"context"

	"github.com/arcbjorn/store-management-system/pb/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore  UserStore
	jwtManager *JWTManager
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{userStore, jwtManager}
}

func (server *AuthServer) Login(
	ctx context.Context,
	req *auth.LoginRequest,
) (*auth.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: &v", err)
	}
	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username / password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot genmerate access token")
	}

	res := &auth.LoginResponse{AccessToken: token}
	return res, nil
}
