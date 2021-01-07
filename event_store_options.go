package edatpgx

import (
	"github.com/stackus/edat/log"
)

type EventStoreOption func(*EventStore)

func WithEventStoreTableName(tableName string) EventStoreOption {
	return func(store *EventStore) {
		store.tableName = tableName
	}
}

func WithEventStoreLogger(logger log.Logger) EventStoreOption {
	return func(store *EventStore) {
		store.logger = logger
	}
}
