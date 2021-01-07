package edatpgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

func ReceiverSessionMiddleware(conn *pgxpool.Pool, logger log.Logger) func(msg.MessageReceiver) msg.MessageReceiver {
	return func(next msg.MessageReceiver) msg.MessageReceiver {
		return msg.ReceiveMessageFunc(func(ctx context.Context, message msg.Message) (err error) {
			var tx pgx.Tx

			tx, err = conn.Begin(ctx)
			if err != nil {
				logger.Error("error while starting the request transaction", log.Error(err))
				return fmt.Errorf("failed to start transaction: %s", err.Error())
			}

			txCtx := context.WithValue(ctx, pgxTxKey, tx)

			defer func() {
				p := recover()
				switch {
				case p != nil:
					txErr := tx.Rollback(ctx)
					if txErr != nil {
						logger.Error("error while rolling back the message receiver transaction during panic", log.Error(txErr))
					}
					panic(p)
				case err != nil:
					txErr := tx.Rollback(ctx)
					if txErr != nil {
						logger.Error("error while rolling back the message receiver transaction", log.Error(txErr))
					}
				default:
					txErr := tx.Commit(ctx)
					if txErr != nil {
						logger.Error("error while committing the message receiver transaction", log.Error(txErr))
					}
				}
			}()

			err = next.ReceiveMessage(txCtx, message)

			return err
		})
	}
}
