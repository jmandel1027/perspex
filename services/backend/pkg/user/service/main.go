package service

import (
	"context"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
	"github.com/jmandel1027/perspex/services/backend/pkg/user/repository"
)

// IUserService for interacting with Users
type IUserService interface {
	RegisterUser(ctx context.Context, in *users.RegisterUserRequest) (*users.User, error)
	RetrieveUser(ctx context.Context, in *users.RetrieveUserRequest) (*users.User, error)
}

// UserService structs
type UserService struct {
	mu   *sync.RWMutex
	repo *repository.UserRepository
	users.UnimplementedUserServiceServer
}

// NewUserService for connecting to the repository
func NewUserService() *UserService {
	repo := repository.NewUserRepository()
	service := &UserService{
		mu:   &sync.RWMutex{},
		repo: repo,
	}

	return service
}

// RegisterUser by RegisterUserRequest
func (svc *UserService) RegisterUser(ctx context.Context, in *users.RegisterUserRequest) (*users.User, error) {
	return &users.User{}, nil
}

// RetrieveUser fetches a user by ID
func (svc *UserService) RetrieveUser(ctx context.Context, in *users.RetrieveUserRequest) (*users.User, error) {
	otelzap.Ctx(ctx).Info("RetrieveUser: ", zap.Int64("id", in.Id))

	record, err := svc.repo.FindUserById(ctx, in.Id)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error retrieving user: ", zap.Error(err))
		return nil, err
	}

	user := &users.User{
		Id:        record.ID,
		AuthId:    "",
		Email:     record.Email,
		FirstName: record.FirstName,
		// I KNOW THIS IS WRONG
		LastName: record.FullName,
	}

	return user, nil
}
