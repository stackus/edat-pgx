package edatpgx

import (
	"github.com/stackus/edat/log"
)

type SagaInstanceStoreOption func(*SagaInstanceStore)

func WithSagaInstanceStoreTableName(tableName string) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.tableName = tableName
	}
}

func WithSagaInstanceStoreLogger(logger log.Logger) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.logger = logger
	}
}
