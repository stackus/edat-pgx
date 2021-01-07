package edatpgx

import (
	"github.com/stackus/edat/log"
)

type MessageStoreOption func(*MessageStore)

func WithMessageStoreTableName(tableName string) MessageStoreOption {
	return func(store *MessageStore) {
		store.tableName = tableName
	}
}

func WithMessageStoreLogger(logger log.Logger) MessageStoreOption {
	return func(store *MessageStore) {
		store.logger = logger
	}
}
