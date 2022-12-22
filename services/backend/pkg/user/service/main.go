package service

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmandel1027/perspex/schemas/perspex/pkg/models"
	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
	"github.com/jmandel1027/perspex/services/backend/pkg/user/repository"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// IUserService for interacting with Users
type IUserService interface {
	DeleteUser(context.Context, *users.UserInputRequest) (*users.User, error)
	ModifyUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error)
	RegisterUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error)
	RetrieveUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error)
	RetrieveUsers(ctx context.Context, in *users.UsersByIdRequest) (*users.Users, error)
	RetrieveUsersPage(ctx context.Context, in *users.UsersPageRequest) (*users.UsersPage, error)
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
	svc.mu.RLock()

	u := &models.User{
		ID:        in.User.Id,
		Email:     in.User.Email,
		FirstName: in.User.FirstName,
		LastName:  in.User.LastName,
	}

	record, err := svc.repo.UpdateUser(ctx, u)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error modifying user: ", zap.Error(err))
		return nil, err
	}

	user := &users.User{
		Id:        record.ID,
		AuthId:    "",
		Email:     record.Email,
		FirstName: record.FirstName,
		LastName:  record.LastName,
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}

	defer svc.mu.RUnlock()

	return user, nil
}

// RegisterUser by RegisterUserRequest
func (svc *UserService) RegisterUser(ctx context.Context, in *users.UserInputRequest) (*users.User, error) {
	svc.mu.RLock()

	u := &models.User{
		Email:     in.User.Email,
		FirstName: in.User.FirstName,
		LastName:  in.User.LastName,
	}

	record, err := svc.repo.CreateUser(ctx, u)
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
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}

	defer svc.mu.RUnlock()

	return user, nil
}

// RetrieveUser fetches a user by ID
func (svc *UserService) RetrieveUser(ctx context.Context, in *users.UserByIdRequest) (*users.User, error) {
	svc.mu.RLock()

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

func (svc *UserService) RetrieveUsers(ctx context.Context, in *users.UsersByIdRequest) (*users.Users, error) {
	svc.mu.RLock()

	records, err := svc.repo.FindUsersByIds(ctx, in.Ids)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error retrieving users: ", zap.Error(err))
		return nil, err
	}

	usr := make([]*users.User, len(records))

	for i, record := range records {

		user := &users.User{
			Id:        record.ID,
			AuthId:    "",
			Email:     record.Email,
			FirstName: record.FirstName,
			LastName:  record.LastName,
			CreatedAt: timestamppb.New(record.CreatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		}

		usr = append(usr[:i], user)
	}

	res := &users.Users{
		Users: usr,
	}

	defer svc.mu.RUnlock()

	return res, nil
}

func (svc *UserService) RetrieveUsersPage(ctx context.Context, in *users.UsersPageRequest) (*users.UsersPage, error) {
	return &users.UsersPage{}, nil
}
