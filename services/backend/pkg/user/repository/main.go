package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmandel1027/perspex/schemas/perspex/pkg/models"
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// UserRepository --
type UserRepository struct {
	cfg *config.BackendConfig
}

// IUserRepository is interface for MaterialRepository
type IUserRepository interface {
	CreateUser(ctx context.Context, record *models.User) (res *models.User, err error)
	FindUserById(ctx context.Context, id int64) (res *models.User, err error)
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
	panic("not implemented")
}

// FindUserById register's a new user
func (repo *UserRepository) FindUserById(ctx context.Context, id int64) (res *models.User, err error) {
	err = postgres.InTx(ctx, sql.TxOptions{ReadOnly: true}, func(tx *postgres.Tx) error {
		res, err = models.Users(models.UserWhere.ID.EQ(id)).One(ctx, tx)
		if err != nil {
			warning := fmt.Sprintf("Couldn't retrieve user: %s", err)
			otelzap.L().Ctx(ctx).Error(warning)
			return errors.New(warning)
		}

		return nil
	})

	return
}
