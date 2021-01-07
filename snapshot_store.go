package edatpgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/edat/log"
)

type SnapshotStore struct {
	tableName string
	client    Client
	strategy  es.SnapshotStrategy
	next      es.AggregateRootStore
	logger    log.Logger
}

var _ es.AggregateRootStore = (*SnapshotStore)(nil)

func NewSnapshotStore(client Client, options ...SnapshotStoreOption) es.AggregateRootStoreMiddleware {
	s := &SnapshotStore{
		tableName: DefaultSnapshotTableName,
		client:    client,
		strategy:  es.DefaultSnapshotStrategy,
	}

	for _, option := range options {
		option(s)
	}

	return func(next es.AggregateRootStore) es.AggregateRootStore {
		s.next = next
		return s
	}
}

func (s *SnapshotStore) Load(ctx context.Context, root *es.AggregateRoot) error {
	var snapshotName string
	var data []byte
	var version int

	name := root.AggregateName()
	id := root.AggregateID()

	row := s.client.QueryRow(ctx, fmt.Sprintf(loadSnapshotSQL, s.tableName), name, id)
	err := row.Scan(&snapshotName, &data, &version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s.next.Load(ctx, root)
		}
		return err
	}

	snapshot, err := core.DeserializeSnapshot(snapshotName, data)
	if err != nil {
		return err
	}

	err = root.LoadSnapshot(snapshot, version)
	if err != nil {
		return err
	}

	return s.next.Load(ctx, root)
}

func (s *SnapshotStore) Save(ctx context.Context, root *es.AggregateRoot) error {
	err := s.next.Save(ctx, root)
	if err != nil {
		return err
	}

	if !s.strategy.ShouldSnapshot(root) {
		return nil
	}

	snapshot, err := root.Aggregate().ToSnapshot()
	if err != nil {
		return err
	}

	data, err := core.SerializeSnapshot(snapshot)
	if err != nil {
		return err
	}

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()

	_, err = s.client.Exec(ctx, fmt.Sprintf(saveSnapshotSQL, s.tableName), name, id, snapshot.SnapshotName(), data, version)
	if err != nil {
		return err
	}

	return nil
}
