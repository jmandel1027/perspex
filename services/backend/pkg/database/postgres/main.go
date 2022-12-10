package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	// import the `pgx` driver for use in `sql.Open`
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

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
)

// TxFunc is a function passed to a transaction block.
type TxFunc func(tx *Tx) error

// key type is unexported to prevent collisions with context keys defined elsewhere.
type key int

// txKey is the context key for the request-scoped transaction, arbitrarily set to 42.
const txWriteKey key = 42

// txKey is the context key for the request-scoped transaction, arbitrarily set to 42.
const txReaderKey key = 43

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
func InTx(ctx context.Context, opts sql.TxOptions, f func(tx *Tx) error) error {
	if ctx.Err() != nil {
		return ErrTXRequiresActiveCtx
	}

	key, err := WhichConnection(ctx, opts)
	if err != nil {
		return ErrTXRequiresActiveCtx
	}

	tx, ok := FromContext(ctx, *key)
	if !ok {
		return ErrTXRequiresActiveCtx
	}

	tx.Lock()
	defer tx.Unlock()
	return f(tx)
}

// WhichConnection returns the key for the connection to use for the given transaction options.
func WhichConnection(ctx context.Context, opts sql.TxOptions) (*key, error) {
	if ctx.Err() != nil {
		return nil, ErrTXRequiresActiveCtx
	}

	var key key
	if opts.ReadOnly {
		key = txReaderKey
	} else {
		key = txWriteKey
	}

	return &key, nil
}

// Open opens a database connection to both writer and reader.
func Open(cfg config.BackendConfig) (*DB, error) {
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
	otelzap.L().Ctx(ctx).Info("BeginTx")
	if ctx.Err() != nil {
		return nil, ErrTXRequiresActiveCtx
	}

	if opts == nil {
		return nil, ErrTXRequiresActiveCtx
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		otelzap.L().Ctx(ctx).Error("BeginTx Error: ", zap.Error(err))
		return nil, err
	}

	otelzap.L().Ctx(ctx).Info("BeginTx Success")
	return &Tx{Tx: tx}, nil
}

// Execute runs a tx-scoped function, commiting on success and rolling back on failure.
func (tx *Tx) Execute(fn TxFunc) (err error) {

	defer func() {
		if p := recover(); err != nil || p != nil {
			tx.Rollback()
			otelzap.L().Info("rolling back")
		} else {
			tx.Commit()
			otelzap.L().Info("committing")
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
func WithTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn TxFunc) error {
	tx, err := BeginTx(ctx, db, opts)
	if err != nil {
		return err
	}

	return tx.Execute(fn)
}

// FromContext extracts an active Postgres transaction from a context.
func FromContext(ctx context.Context, key key) (*Tx, bool) {
	tx, ok := ctx.Value(key).(*Tx)
	return tx, ok
}

// NewContext returns a new context that carries a provided Postgres transaction.
func NewContext(ctx context.Context, key key, tx *Tx) context.Context {
	return context.WithValue(ctx, key, tx)
}
