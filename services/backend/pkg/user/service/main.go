package service

import (
	"context"
	"fmt"
	"sync"

	connect "github.com/bufbuild/connect-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmandel1027/perspex/schemas/perspex/pkg/models"
	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"

	"github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1/usersconnect"
	"github.com/jmandel1027/perspex/services/backend/pkg/user/repository"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// IUserService for interacting with Users
type IUserService interface {
	DeleteUser(ctx context.Context, rec *connect.Request[users.DeleteUserRequest]) (*connect.Response[users.DeleteUserResponse], error)
	ModifyUser(ctx context.Context, rec *connect.Request[users.ModifyUserRequest]) (*connect.Response[users.ModifyUserResponse], error)
	RegisterUser(ctx context.Context, rec *connect.Request[users.RegisterUserRequest]) (*connect.Response[users.RegisterUserRequest], error)
	RetrieveUser(ctx context.Context, rec *connect.Request[users.RetrieveUserRequest]) (*connect.Response[users.RetrieveUserResponse], error)
	RetrieveUsers(ctx context.Context, rec *connect.Request[users.RetrieveUsersRequest]) (*connect.Response[users.RetrieveUsersResponse], error)
	RetrieveUsersPage(ctx context.Context, rec *connect.Request[users.RetrieveUsersPageRequest]) (*connect.Response[users.RetrieveUsersPageResponse], error)
}

// UserService structs
type UserService struct {
	mu   *sync.RWMutex
	repo *repository.UserRepository
	usersconnect.UnimplementedUserServiceHandler
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

func (svc *UserService) DeleteUser(ctx context.Context, rec *connect.Request[users.DeleteUserRequest]) (*connect.Response[users.DeleteUserResponse], error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

// ModidifyUser
func (svc *UserService) ModifyUser(ctx context.Context, rec *connect.Request[users.ModifyUserRequest]) (*connect.Response[users.ModifyUserResponse], error) {
	svc.mu.RLock()

	u := &models.User{
		ID:        rec.Msg.User.Id,
		Email:     rec.Msg.User.Email,
		FirstName: rec.Msg.User.FirstName,
		LastName:  rec.Msg.User.LastName,
	}

	record, err := svc.repo.UpdateUser(ctx, u)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error modifying user: ", zap.Error(err))
		return nil, err
	}

	res := connect.NewResponse(&users.ModifyUserResponse{
		User: &users.User{
			Id:        record.ID,
			AuthId:    "",
			Email:     record.Email,
			FirstName: record.FirstName,
			LastName:  record.LastName,
			CreatedAt: timestamppb.New(record.CreatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		},
	})

	defer svc.mu.RUnlock()

	return res, nil
}

// RegisterUser by RegisterUserRequest
func (svc *UserService) RegisterUser(ctx context.Context, rec *connect.Request[users.RegisterUserRequest]) (*connect.Response[users.RegisterUserResponse], error) {
	svc.mu.RLock()

	u := &models.User{
		Email:     rec.Msg.User.Email,
		FirstName: rec.Msg.User.FirstName,
		LastName:  rec.Msg.User.LastName,
	}

	record, err := svc.repo.CreateUser(ctx, u)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error retrieving user: ", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&users.RegisterUserResponse{
		User: &users.User{
			Id:        record.ID,
			AuthId:    "",
			Email:     record.Email,
			FirstName: record.FirstName,
			LastName:  record.LastName,
			CreatedAt: timestamppb.New(record.CreatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		},
	})

	defer svc.mu.RUnlock()

	return res, nil
}

// RetrieveUser fetches a user by ID
func (svc *UserService) RetrieveUser(ctx context.Context, rec *connect.Request[users.RetrieveUserRequest]) (*connect.Response[users.RetrieveUserResponse], error) {
	svc.mu.RLock()

	record, err := svc.repo.FindUserById(ctx, rec.Msg.Id)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error retrieving user", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if record == nil {
		otelzap.Ctx(ctx).Error("Error retrieving user", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	res := connect.NewResponse(&users.RetrieveUserResponse{
		User: &users.User{
			Id:        record.ID,
			AuthId:    "",
			Email:     record.Email,
			FirstName: record.FirstName,
			LastName:  record.LastName,
			CreatedAt: timestamppb.New(record.CreatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		},
	})

	defer svc.mu.RUnlock()

	return res, nil
}

func (svc *UserService) RetrieveUsers(ctx context.Context, rec *connect.Request[users.RetrieveUsersRequest]) (*connect.Response[users.RetrieveUsersResponse], error) {
	svc.mu.RLock()

	records, err := svc.repo.FindUsersByIds(ctx, rec.Msg.Ids)
	if err != nil {
		otelzap.Ctx(ctx).Error("Error retrieving users: ", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, err)
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

	res := connect.NewResponse(&users.RetrieveUsersResponse{
		Users: usr,
	})

	defer svc.mu.RUnlock()

	return res, nil
}

func (svc *UserService) RetrieveUsersPage(ctx context.Context, rec *connect.Request[users.RetrieveUsersPageRequest]) (*connect.Response[users.RetrieveUsersPageResponse], error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveUsersPage not implemented")
}
