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
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
)

// UserRepository --
type UserRepository struct {
	cfg *config.BackendConfig
}

// IUserRepository is interface for MaterialRepository
type IUserRepository interface {
	CreateUser(ctx context.Context, record *models.User) (res *models.User, err error)
	FindUserById(ctx context.Context, id int64) (res *models.User, err error)
	FindUsersByIds(ctx context.Context, ids []int64) (res []*models.User, err error)
	UpdateUser(ctx context.Context, record *models.User) (res *models.User, err error)
}

// NewUserRepository Creates a new Material repo instance
func NewUserRepository() *UserRepository {
	cfg, _ := config.New()
	return &UserRepository{
		&cfg,
	}
}

// CreateUser register's a new user
func (repo *UserRepository) CreateUser(ctx context.Context, record *models.User) (res *models.User, err error) {
	err = postgres.InTx(ctx, postgres.StdTxOpts, func(tx *postgres.Tx) error {
		if err = record.Insert(ctx, tx, boil.Infer()); err != nil {
			warning := fmt.Sprintf("Couldn't register user: %s", err)
			otelzap.L().Ctx(ctx).Error(warning)
			return errors.New(warning)
		}

		res = record
		return nil
	})

	return
}

// FindUserById register's a new user
func (repo *UserRepository) FindUserById(ctx context.Context, id int64) (res *models.User, err error) {
	err = postgres.InTx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *postgres.Tx) error {
		otelzap.L().Ctx(ctx).Info("attempting to fetch user")
		res, err = models.Users(models.UserWhere.ID.EQ(id)).One(ctx, tx)
		if err != nil && err == sql.ErrNoRows {
			otelzap.L().Ctx(ctx).Info("no user found")
			return nil
		}

		if err != nil {
			warning := fmt.Sprintf("Couldn't retrieve user: %s", err)
			otelzap.L().Ctx(ctx).Error(warning)
			return errors.New(warning)
		}

		return nil
	})

	otelzap.L().Ctx(ctx).Info("done attempting to find user")

	return
}

// FindUsersByIds finds users by ids
func (repo *UserRepository) FindUsersByIds(ctx context.Context, ids []int64) (res []*models.User, err error) {
	err = postgres.InTx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *postgres.Tx) error {
		res, err = models.Users(qm.Where("id = ANY ($1)", ids)).All(ctx, tx)
		if err != nil && err == sql.ErrNoRows {
			return nil
		}

		if err != nil {
			warning := fmt.Sprintf("Couldn't retrieve users: %s", err)
			otelzap.L().Ctx(ctx).Error(warning)
			return errors.New(warning)
		}

		return nil
	})

	return
}

// UpdateUser modifies a user
func (repo *UserRepository) UpdateUser(ctx context.Context, record *models.User) (res *models.User, err error) {
	err = postgres.InTx(ctx, postgres.StdTxOpts, func(tx *postgres.Tx) error {
		if _, err = record.Update(ctx, tx, boil.Infer()); err != nil {
			warning := fmt.Sprintf("Couldn't update user: %s", err)
			otelzap.L().Ctx(ctx).Error(warning)
			return errors.New(warning)
		}

		res = record
		return nil
	})

	return
}
