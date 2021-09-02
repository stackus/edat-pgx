package edatpgx

import (
	"github.com/stackus/edat/log"
)

type SessionMiddlewareOption interface {
	configureSessionMiddleware(*SessionMiddlewareConfig)
}

type SessionMiddlewareConfig struct {
	logger log.Logger
}

func NewSessionMiddlewareConfig() *SessionMiddlewareConfig {
	return &SessionMiddlewareConfig{
		logger: log.NewNopLogger(),
	}
}
