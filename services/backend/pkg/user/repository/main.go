package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/jmandel1027/perspex/schemas/perspex/pkg/models"
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
)

// UserRepository --
type UserRepository struct {
	cfg *config.BackendConfig
}

// IUserRepository is interface for MaterialRepository
type IUserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, record *models.User) (*models.User, error)
	FindUserById(ctx context.Context, tx *sql.Tx, id int64) (*models.User, error)
	FindUsersByIds(ctx context.Context, tx *sql.Tx, ids []int64) ([]*models.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, record *models.User) (*models.User, error)
}

// NewUserRepository Creates a new Material repo instance
func NewUserRepository() *UserRepository {
	cfg, _ := config.New()
	return &UserRepository{
		&cfg,
	}
}

// CreateUser register's a new user
func (repo *UserRepository) CreateUser(ctx context.Context, tx *sql.Tx, record *models.User) (*models.User, error) {

	if err := record.Insert(ctx, tx, boil.Infer()); err != nil {
		warning := fmt.Sprintf("Couldn't register user: %s", err)
		otelzap.L().Ctx(ctx).Error(warning)
		return nil, errors.New(warning)
	}

	return record, nil
}

// FindUserById register's a new user
func (repo *UserRepository) FindUserById(ctx context.Context, tx *sql.Tx, id int64) (*models.User, error) {

	otelzap.L().Ctx(ctx).Info("attempting to fetch user")
	record, err := models.Users(models.UserWhere.ID.EQ(id)).One(ctx, tx)
	if err != nil && err == sql.ErrNoRows {
		otelzap.L().Ctx(ctx).Info("no user found")
		return nil, nil
	}

	if err != nil {
		warning := fmt.Sprintf("Couldn't retrieve user: %s", err)
		otelzap.L().Ctx(ctx).Error(warning)
		return nil, errors.New(warning)
	}

	otelzap.L().Ctx(ctx).Info("done attempting to find user")

	return record, nil
}

// FindUsersByIds finds users by ids
func (repo *UserRepository) FindUsersByIds(ctx context.Context, tx *sql.Tx, ids []int64) ([]*models.User, error) {
	record, err := models.Users(qm.Where("id = ANY ($1)", ids)).All(ctx, tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		warning := fmt.Sprintf("Couldn't retrieve users: %s", err)
		otelzap.L().Ctx(ctx).Error(warning)
		return nil, errors.New(warning)
	}

	return record, nil
}

// UpdateUser modifies a user
func (repo *UserRepository) UpdateUser(ctx context.Context, tx *sql.Tx, record *models.User) (*models.User, error) {

	if _, err := record.Update(ctx, tx, boil.Infer()); err != nil {
		warning := fmt.Sprintf("Couldn't update user: %s", err)
		otelzap.L().Ctx(ctx).Error(warning)
		return nil, errors.New(warning)
	}

	return record, nil
}
