package edatpgx

import (
	"github.com/stackus/edat/es"
	"github.com/stackus/edat/log"
)

type SnapshotStoreOption func(*SnapshotStore)

func WithSnapshotStoreTableName(tableName string) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.tableName = tableName
	}
}

func WithSnapshotStoreStrategy(strategy es.SnapshotStrategy) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.strategy = strategy
	}
}

func WithSnapshotStoreLogger(logger log.Logger) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.logger = logger
	}
}
