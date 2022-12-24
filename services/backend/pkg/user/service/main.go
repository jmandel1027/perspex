package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	connect "github.com/bufbuild/connect-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmandel1027/perspex/schemas/perspex/pkg/models"
	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"

	"github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1/usersconnect"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	"github.com/jmandel1027/perspex/services/backend/pkg/user/repository"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// IUserService for interacting with Users
type IUserService interface {
	DeleteUser(ctx context.Context, rec *connect.Request[users.DeleteUserRequest]) (*connect.Response[users.DeleteUserResponse], error)
	ModifyUser(ctx context.Context, rec *connect.Request[users.ModifyUserRequest]) (*connect.Response[users.ModifyUserResponse], error)
	RegisterUser(ctx context.Context, rec *connect.Request[users.RegisterUserRequest]) (*connect.Response[users.RegisterUserResponse], error)
	RetrieveUser(ctx context.Context, rec *connect.Request[users.RetrieveUserRequest]) (*connect.Response[users.RetrieveUserResponse], error)
	RetrieveUsers(ctx context.Context, rec *connect.Request[users.RetrieveUsersRequest]) (*connect.Response[users.RetrieveUsersResponse], error)
	RetrieveUsersPage(ctx context.Context, rec *connect.Request[users.RetrieveUsersPageRequest]) (*connect.Response[users.RetrieveUsersPageResponse], error)
}

// UserService structs
type UserService struct {
	mu   *sync.RWMutex
	db   *postgres.DB
	repo *repository.UserRepository
	usersconnect.UnimplementedUserServiceHandler
}

// NewUserService for connecting to the repository
func NewUserService(dbs *postgres.DB) *UserService {
	repo := repository.NewUserRepository()
	service := &UserService{
		mu:   &sync.RWMutex{},
		db:   dbs,
		repo: repo,
	}

	return service
}

func (svc *UserService) DeleteUser(ctx context.Context, rec *connect.Request[users.DeleteUserRequest]) (*connect.Response[users.DeleteUserResponse], error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

// ModidifyUser
func (svc *UserService) ModifyUser(ctx context.Context, rec *connect.Request[users.ModifyUserRequest]) (*connect.Response[users.ModifyUserResponse], error) {
	var record *models.User
	svc.mu.RLock()

	u := &models.User{
		ID:        rec.Msg.User.Id,
		Email:     rec.Msg.User.Email,
		FirstName: rec.Msg.User.FirstName,
		LastName:  rec.Msg.User.LastName,
	}

	postgres.WithTx(ctx, svc.db.Writer, postgres.StdTxOpts, func(tx *sql.Tx) (err error) {
		record, err = svc.repo.UpdateUser(ctx, tx, u)
		if err != nil {
			otelzap.Ctx(ctx).Error("Error modifying user: ", zap.Error(err))
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})

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
	var record *models.User
	svc.mu.RLock()

	u := &models.User{
		Email:     rec.Msg.User.Email,
		FirstName: rec.Msg.User.FirstName,
		LastName:  rec.Msg.User.LastName,
	}

	postgres.WithTx(ctx, svc.db.Reader, postgres.ReadOnlyTxOpts, func(tx *sql.Tx) (err error) {
		record, err = svc.repo.CreateUser(ctx, tx, u)
		if err != nil {
			otelzap.Ctx(ctx).Error("Error retrieving user: ", zap.Error(err))
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})

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
	var record *models.User
	svc.mu.RLock()

	postgres.WithTx(ctx, svc.db.Reader, postgres.ReadOnlyTxOpts, func(tx *sql.Tx) (err error) {
		record, err = svc.repo.FindUserById(ctx, tx, rec.Msg.Id)
		if err != nil {
			otelzap.Ctx(ctx).Error("Error retrieving user", zap.Error(err))
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})

	if record == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	res := connect.NewResponse(&users.RetrieveUserResponse{
		User: &users.User{
			Id:        record.ID,
			AuthId:    "",
			Email:     record.Email,
			FirstName: record.FirstName,
			LastName:  record.LastName,
			//CreatedAt: timestamppb.New(time.Now()),
			//UpdatedAt: timestamppb.New(time.Now()),
			//CreatedAt: timestamppb.New(record.CreatedAt),
			// UpdatedAt: timestamppb.New(record.UpdatedAt),
		},
	})

	defer svc.mu.RUnlock()

	return res, nil
}

func (svc *UserService) RetrieveUsers(ctx context.Context, rec *connect.Request[users.RetrieveUsersRequest]) (*connect.Response[users.RetrieveUsersResponse], error) {
	var records []*models.User
	svc.mu.RLock()

	postgres.WithTx(ctx, svc.db.Reader, postgres.ReadOnlyTxOpts, func(tx *sql.Tx) (err error) {
		records, err = svc.repo.FindUsersByIds(ctx, tx, rec.Msg.Ids)
		if err != nil {
			otelzap.Ctx(ctx).Error("Error retrieving users: ", zap.Error(err))
			return connect.NewError(connect.CodeInternal, err)
		}

		return nil
	})

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
