package service

import (
	"context"
	"sync"

	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
)

// IUserService for interacting with Users
type IUserService interface {
	RegisterUser(ctx context.Context, in *users.RegisterUserRequest) (*users.User, error)
	RetrieveUser(ctx context.Context, in *users.RetrieveUserRequest) (*users.User, error)
}

// UserService structs
type UserService struct {
	mu *sync.RWMutex
	users.UnimplementedUserServiceServer
}

// NewUserService for connecting to the repository
func NewUserService() *UserService {
	service := &UserService{
		mu: &sync.RWMutex{},
	}

	return service
}

// RegisterUser by RegisterUserRequest
func (svc *UserService) RegisterUser(ctx context.Context, in *users.RegisterUserRequest) (*users.User, error) {
	return &users.User{}, nil
}

// RetrieveUser fetches a user by ID
func (svc *UserService) RetrieveUser(ctx context.Context, in *users.RetrieveUserRequest) (*users.User, error) {
	return &users.User{}, nil
}
