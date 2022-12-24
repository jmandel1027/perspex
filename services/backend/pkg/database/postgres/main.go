package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	// import the `pgx` driver for use in `sql.Open`
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"

	"github.com/jmandel1027/perspex/services/backend/pkg/config"
)

// DB is a PostgreSQL connection. It uses the `pgx` driver wrapped by the `pgx/stdlib` compatibility layer.
type DB struct {
	Writer *sql.DB
	Reader *sql.DB
}

// Tx is a PostgreSQL transaction.
type Tx struct {
	*sql.Tx
	mu sync.Mutex
}

// Error strings
var (
	ErrTXRequiresUser      = errors.New("user must be present in context to begin transaction")
	ErrTXRequiresActiveCtx = errors.New("context must not be cancelled to begin transaction")
	ErrTXRequiresOpts      = errors.New("opts must be passed to begin transaction")
)

type key string

const Key key = "postgres"
const Writer key = "writer"
const Reader key = "reader"

// TxFunc is a function passed to a transaction block.
type TxFunc func(tx *sql.Tx) error

// StdTxOpts are standard repeatable read transaction options used for most database operations.
var StdTxOpts = &sql.TxOptions{Isolation: sql.LevelRepeatableRead}

// CommittedTxOpts are standard repeatable read transaction options used for most database operations.
var CommittedTxOpts = &sql.TxOptions{Isolation: sql.LevelReadCommitted}

// ReadOnlyTxOpts are read-only transactions options used for long-lasting read-only transactions that need to
// read data committed by other concurrent transactions.
var ReadOnlyTxOpts = &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}

// Open opens a database connection to both writer and reader.
func Open(cfg *config.BackendConfig) (*DB, error) {
	writer, err := sql.Open("pgx", cfg.WriterPG.GetDataSourceName())
	if err != nil {
		otelzap.L().Error("WriterPG Error: ", zap.Error(err))
		return nil, err
	}

	writer.SetMaxOpenConns(cfg.WriterPG.GetMaxOpenConns())
	writer.SetMaxIdleConns(cfg.WriterPG.GetMaxIdleConns())
	writer.SetConnMaxLifetime(cfg.WriterPG.GetConnMaxLifetime())

	reader, err := sql.Open("pgx", cfg.ReaderPG.GetDataSourceName())
	if err != nil {
		otelzap.L().Error("ReaderPG Error: ", zap.Error(err))
		return nil, err
	}

	reader.SetMaxOpenConns(cfg.ReaderPG.GetMaxOpenConns())
	reader.SetMaxIdleConns(cfg.ReaderPG.GetMaxIdleConns())
	reader.SetConnMaxLifetime(cfg.ReaderPG.GetConnMaxLifetime())

	return &DB{writer, reader}, nil
}

// BeginTx initializes a transaction.
func BeginTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions) (*Tx, error) {
	otelzap.L().Info("attempting to begin transaction")

	if ctx.Err() != nil {
		otelzap.L().Info("rolling back")
		return nil, ErrTXRequiresActiveCtx
	}

	if opts == nil {
		otelzap.L().Info("reequired opts")
		return nil, ErrTXRequiresOpts
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}

	otelzap.L().Info("instantiated transaction")

	return &Tx{Tx: tx}, nil
}

// Execute runs a tx-scoped function, commiting on success and rolling back on failure.
func (tx *Tx) Execute(fn TxFunc) (err error) {
	tx.Lock()

	if err := fn(tx.Tx); err != nil {
		return err
	}

	defer func() {
		if p := recover(); err != nil || p != nil {
			if err := tx.Tx.Rollback(); err != nil {
				otelzap.L().Ctx(context.Background()).Error("Failed to rollback transaction", zap.Error(err))
			}
			otelzap.L().Info("rolling back")
		} else {
			if err := tx.Commit(); err != nil {
				otelzap.L().Ctx(context.Background()).Error("Failed to commit transaction", zap.Error(err))
			}
			otelzap.L().Info("committing")
		}

	}()

	defer tx.Unlock()

	return nil
}

// Lock locks the transaction, preventing concurrent use.
func (tx *Tx) Lock() {
	tx.mu.Lock()
}

// Unlock unlocks the transaction, enabling concurrent use.
func (tx *Tx) Unlock() {
	tx.mu.Unlock()
}

// WithTx creates a transaction block, commits on success, and rolls back on failure.
func WithTx(ctx context.Context, opts *sql.TxOptions, fn TxFunc) error {

	cfg, err := config.New()
	if err != nil {
		return err
	}

	db, ok := FromContext(ctx, Key)
	if !ok {
		otelzap.Ctx(ctx).Error("Transaction requires active ctx")
		return ErrTXRequiresActiveCtx
	}

	var conn *sql.DB
	if opts.ReadOnly {
		conn = db.Reader
	} else {
		conn = db.Writer
	}

	if cfg.Log.Verbose {
		writer := &zapio.Writer{Log: otelzap.Ctx(ctx).Logger().Logger, Level: zap.DebugLevel}

		defer writer.Close()

		boil.DebugMode = true
		boil.DebugWriter = writer
	}

	tx, err := BeginTx(ctx, conn, opts)
	if err != nil {
		return err
	}

	return tx.Execute(fn)
}

// FromContext extracts an active Postgres transaction from a context.
func FromContext(ctx context.Context, key key) (*DB, bool) {
	db, ok := ctx.Value(key).(*DB)
	return db, ok
}

// NewContext returns a new context that carries a provided Postgres transaction.
func NewContext(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, Key, db)
}
