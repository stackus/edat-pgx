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

type EventStore struct {
	tableName string
	client    Client
	logger    log.Logger
}

var _ es.AggregateRootStore = (*EventStore)(nil)

func NewEventStore(client Client, options ...EventStoreOption) *EventStore {
	store := &EventStore{
		tableName: DefaultEventTableName,
		client:    client,
		logger:    log.DefaultLogger,
	}

	for _, option := range options {
		option(store)
	}

	return store
}

func (s *EventStore) Load(ctx context.Context, root *es.AggregateRoot) error {
	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()

	rows, err := s.client.Query(ctx, fmt.Sprintf(loadEventsSQL, s.tableName), name, id, version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var event core.Event
		var eventName string
		var data []byte

		err = rows.Scan(&eventName, &data)
		if err != nil {
			return err
		}
		event, err = core.DeserializeEvent(eventName, data)
		if err != nil {
			return err
		}
		err = root.LoadEvent(event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EventStore) Save(ctx context.Context, root *es.AggregateRoot) (err error) {
	var tx pgx.Tx

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.Version()

	tx, err = s.client.Begin(ctx)
	if err != nil {
		return
	}

	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			_ = tx.Rollback(ctx)
		default:
			err = tx.Commit(ctx)
		}
	}()

	correlationID := core.GetCorrelationID(ctx)
	causationID := core.GetRequestID(ctx)

	for i, event := range root.Events() {
		var data []byte

		data, err = core.SerializeEvent(event)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, fmt.Sprintf(writeEventSQL, s.tableName), name, id, correlationID, causationID, version+i+1, event.EventName(), data)
		if err != nil {
			return err
		}
	}

	return
}
