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

// TxFunc is a function passed to a transaction block.
type TxFunc func(tx *Tx) error

// SqlTxFunc is a function passed to a transaction block.
type SqlTxFunc func(tx *sql.Tx) error

// key type is unexported to prevent collisions with context keys defined elsewhere.
type Key int

// txKey is the context key for the request-scoped transaction, arbitrarily set to 42.
const txWriteKey Key = 42

// txKey is the context key for the request-scoped transaction, arbitrarily set to 42.
const txReaderKey Key = 43

// StdTxOpts are standard repeatable read transaction options used for most database operations.
var StdTxOpts = &sql.TxOptions{Isolation: sql.LevelRepeatableRead}

// CommittedTxOpts are standard repeatable read transaction options used for most database operations.
var CommittedTxOpts = &sql.TxOptions{Isolation: sql.LevelReadCommitted}

// ReadOnlyTxOpts are read-only transactions options used for long-lasting read-only transactions that need to
// read data committed by other concurrent transactions.
var ReadOnlyTxOpts = &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}

// InTx provides the passed function with exclusive access to the Postgres
// transaction within the passed context. A non-nil error will be returned
// if the context is not open or does not contain a transaction.
func InTx(ctx context.Context, opts *sql.TxOptions, f func(tx *Tx) error) error {
	otelzap.Ctx(ctx).Info("in transaction")
	if ctx.Err() != nil {
		otelzap.Ctx(ctx).Info(ctx.Err().Error())
		return ErrTXRequiresActiveCtx
	}

	key, err := WhichConnection(ctx, opts)
	if err != nil {
		otelzap.Ctx(ctx).Error("Requires opts", zap.Error(err))
		return ErrTXRequiresOpts
	}

	tx, ok := FromContext(ctx, *key)
	if !ok {
		otelzap.Ctx(ctx).Error("Transaction requires active ctx")
		return ErrTXRequiresActiveCtx
	}

	tx.Lock()

	defer tx.Unlock()

	return f(tx)
}

// WhichConnection returns the key for the connection to use for the given transaction options.
func WhichConnection(ctx context.Context, opts *sql.TxOptions) (*Key, error) {
	if ctx.Err() != nil {
		return nil, ErrTXRequiresActiveCtx
	}

	if opts == nil {
		return nil, ErrTXRequiresOpts
	}

	var key Key
	if opts.ReadOnly {
		key = txReaderKey
	} else {
		key = txWriteKey
	}

	return &key, nil
}

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
	if ctx.Err() != nil {
		return nil, ErrTXRequiresActiveCtx
	}

	if opts == nil {
		return nil, ErrTXRequiresOpts
	}

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	if cfg.Log.Verbose {
		writer := &zapio.Writer{
			Log:   otelzap.Ctx(ctx).Logger().Logger,
			Level: zap.DebugLevel,
		}

		defer writer.Close()

		boil.DebugMode = true
		boil.DebugWriter = writer
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

// Execute runs a tx-scoped function, commiting on success and rolling back on failure.
func (tx *Tx) Execute(fn TxFunc) (err error) {

	defer func() {
		if p := recover(); err != nil || p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}

	}()

	return fn(tx)
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
func WithTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn SqlTxFunc) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.Log.Verbose {
		writer := &zapio.Writer{
			Log:   otelzap.Ctx(ctx).Logger().Logger,
			Level: zap.DebugLevel,
		}

		defer writer.Close()

		boil.DebugMode = true
		boil.DebugWriter = writer
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			if err := tx.Commit(); err != nil {
				// https://golang.org/pkg/database/sql/#Tx
				// After a call to Commit or Rollback, all operations on the
				// transaction fail with ErrTxDone.
				if err == sql.ErrTxDone {
					otelzap.L().Ctx(ctx).Error("Failed to commit transaction", zap.Error(err))
				}
			}
		}
	}()

	return fn(tx)
}

// FromContext extracts an active Postgres transaction from a context.
func FromContext(ctx context.Context, key Key) (*Tx, bool) {
	tx, ok := ctx.Value(key).(*Tx)
	return tx, ok
}

// NewContext returns a new context that carries a provided Postgres transaction.
func NewContext(ctx context.Context, key Key, tx *Tx) context.Context {
	return context.WithValue(ctx, key, tx)
}
