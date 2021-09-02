package edatpgx

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackus/edat/log"
	"google.golang.org/grpc"
)

func RpcSessionUnaryInterceptor(conn *pgxpool.Pool, options ...SessionMiddlewareOption) grpc.UnaryServerInterceptor {
	cfg := NewSessionMiddlewareConfig()

	for _, option := range options {
		option.configureSessionMiddleware(cfg)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tx, err := conn.Begin(ctx)
		if err != nil {
			cfg.logger.Error("error while starting the request transaction", log.Error(err))
			return
		}

		newCtx := context.WithValue(ctx, pgxTxKey, tx)

		defer func() {
			p := recover()
			switch {
			case p != nil:
				txErr := tx.Rollback(ctx)
				if txErr != nil {
					cfg.logger.Error("error while rolling back the rpc request transaction during panic", log.Error(txErr))
				}
				panic(p)
			case err != nil:
				txErr := tx.Rollback(ctx)
				if txErr != nil {
					cfg.logger.Error("error while rolling back the rpc request transaction", log.Error(txErr))
				}
			default:
				txErr := tx.Commit(ctx)
				if txErr != nil {
					cfg.logger.Error("error while committing the rpc request transaction", log.Error(txErr))
				}
			}
		}()

		return handler(newCtx, req)
	}
}

func RpcSessionStreamInterceptor(_ *pgxpool.Pool, options ...SessionMiddlewareOption) grpc.StreamServerInterceptor {
	cfg := NewSessionMiddlewareConfig()

	for _, option := range options {
		option.configureSessionMiddleware(cfg)
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cfg.logger.Error("outbox pattern not yet implemented for streaming connections")
		return handler(srv, ss)
	}
}
