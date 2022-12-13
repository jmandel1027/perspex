package service

import (
	"context"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func (svc *UserService) DeleteUser(context.Context, *users.UserInputRequest) (*users.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

// ModidifyUser
func (svc *UserService) ModifyUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error) {
	return &users.User{}, nil
}

// RegisterUser by RegisterUserRequest
func (svc *UserService) RegisterUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error) {
	return &users.User{}, nil
}

// RetrieveUser fetches a user by ID
func (svc *UserService) RetrieveUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error) {
	svc.mu.RLock()

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
		LastName:  record.LastName,
		CreatedAt: &timestamppb.Timestamp{},
		UpdatedAt: &timestamppb.Timestamp{},
	}

	defer svc.mu.RUnlock()

	return user, nil
}

func (svc *UserService) RetrieveUsers(ctx context.Context, in *users.RetrieveUsersRequest) (*users.Users, error) {
	return &users.Users{}, nil
}

func (svc *UserService) RetrieveUsersPage(ctx context.Context, in *users.RetrieveUsersPageRequest) (*users.UsersPage, error) {
	return &users.UsersPage{}, nil
}
