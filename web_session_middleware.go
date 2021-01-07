package edatpgx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stackus/edat/log"
)

func WebSessionMiddleware(conn *pgxpool.Pool, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := conn.Begin(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "failed to start transaction: %s", err.Error())
				logger.Error("error while starting the request transaction", log.Error(err))
				return
			}

			ww := newStatusWriter(w)

			ctx := context.WithValue(r.Context(), pgxTxKey, tx)

			defer func() {
				p := recover()
				switch {
				case p != nil:
					txErr := tx.Rollback(r.Context())
					if txErr != nil {
						logger.Error("error while rolling back the web request transaction during panic", log.Error(txErr))
					}
					panic(p)
				case ww.status > 399:
					txErr := tx.Rollback(r.Context())
					if txErr != nil {
						logger.Error("error while rolling back the web request transaction", log.Error(txErr))
					}
				default:
					txErr := tx.Commit(r.Context())
					if txErr != nil {
						logger.Error("error while committing the web request transaction", log.Error(txErr))
					}
				}
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func newStatusWriter(w http.ResponseWriter) *statusWriter {
	return &statusWriter{
		ResponseWriter: w,
	}
}

func (s *statusWriter) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}
