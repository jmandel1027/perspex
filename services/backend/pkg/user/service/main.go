package service

import (
	"context"
	"sync"

	user "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
)

// IUserService for interacting with Users
type IUserService interface {
	Echo(ctx context.Context, in *user.EchoRequest) (*user.EchoResponse, error)
}

// UserService structs
type UserService struct {
	mu *sync.RWMutex
	user.UnimplementedUserServiceServer
}

// NewUserService for connecting to the repository
func NewUserService() *UserService {
	service := &UserService{
		mu: &sync.RWMutex{},
	}

	return service
}

// Echo for testing
func (svc *UserService) Echo(ctx context.Context, in *user.EchoRequest) (*user.EchoResponse, error) {
	return &user.EchoResponse{}, nil
}
