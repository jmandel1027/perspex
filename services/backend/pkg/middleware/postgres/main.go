package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type key int

// GetConnectionKey retrieves the connection Key, if any, from the request context.
func GetConnectionKey(ctx context.Context, k postgres.Key) any {
	return ctx.Value(k)
}

// WithoutIdentity strips the db connection, if any, from the provided context.
func WithoutConnectionKey(ctx context.Context, k key) context.Context {
	return context.WithValue(ctx, k, nil)
}

// Errorf is a convenience function that returns an error coded with
// [connect.CodeUnauthenticated].
func Errorf(template string, args ...any) *connect.Error {
	return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(template, args...))
}

// Request describes a single RPC invocation.
type Request struct {
	Spec   connect.Spec
	Peer   connect.Peer
	Header http.Header
}

// Interceptor is a server-side authentication interceptor. In addition to
// rejecting unauthenticated requests, it can optionally attach an identity to
// context of authenticated requests.
type Interceptor struct {
	connection func(context.Context, *Request) (*sql.DB, *sql.TxOptions, error)
}

// New constructs a new Interceptor using the supplied DB pointer and connection options.
// The authentication function must return an error if the request cannot be attached to a read or write tx.
// The error is typically produced with [Errorf], but any error
// will do.
//
// If requests are successfully connected, the interceptor may
// return the transaction. The identity will be attached to the
// context, so subsequent interceptors and application code may access it with
// [GetConnectionKey].
//
// Transaction functions must be safe to call concurrently.
func New(f func(context.Context, *Request) (*sql.DB, *sql.TxOptions, error)) *Interceptor {
	otelzap.L().Info("instantiating interceptors")
	return &Interceptor{f}
}

// WrapUnary implements connect.Interceptor.
func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		conn, settings, err := i.connection(ctx, &Request{
			Spec:   req.Spec(),
			Peer:   req.Peer(),
			Header: req.Header(),
		})

		if err != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			cancel()
		}

		tx, err := postgres.BeginTx(ctx, conn, settings)
		if err != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			cancel()
		}

		key, err := postgres.WhichConnection(ctx, settings)
		if err != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			cancel()
		}

		defer func() {
			if p := recover(); err != nil || p != nil {
				otelzap.L().Info("recovering, rollingback")
				if err := tx.Rollback(); err != nil {
					otelzap.L().Info(err.Error())
				}
			} else {
				otelzap.L().Info("committing")
				if err := tx.Commit(); err != nil {
					otelzap.L().Info(err.Error())
				}
			}
		}()

		return next(postgres.NewContext(ctx, *key, tx), req)
	}
}

// WrapStreamingClient implements connect.Interceptor with a no-op.
func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler implements connect.Interceptor.
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		stream, opts, err := i.connection(ctx, &Request{
			Spec:   conn.Spec(),
			Peer:   conn.Peer(),
			Header: conn.RequestHeader(),
		})
		if err != nil {
			return err
		}

		key, err := postgres.WhichConnection(ctx, opts)
		if err != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			cancel()
		}

		return next(context.WithValue(ctx, key, stream), conn)
	}
}

func Wrap(ctx context.Context, db *sql.DB, txOpts *sql.TxOptions) (*sql.DB, *sql.TxOptions, error) {
	return db, txOpts, nil
}
