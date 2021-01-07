package edatpgx

import (
	"errors"
)

type contextKey int

const (
	DefaultEventTableName        = "events"
	DefaultMessageTableName      = "messages"
	DefaultSagaInstanceTableName = "saga_instances"
	DefaultSnapshotTableName     = "snapshots"

	CreateEventsTableSQL = `CREATE TABLE %s (
    entity_name    text        NOT NULL,
    entity_id      text        NOT NULL,
    correlation_id text        NOT NULL,
    causation_id   text        NOT NULL,
    event_version  int         NOT NULL,
    event_name     text        NOT NULL,
    event_data     bytea       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id, event_version)
)`

	CreateMessagesTableSQL = `CREATE TABLE messages (
    message_id  text        NOT NULL,
    destination text        NOT NULL,
    payload     bytea       NOT NULL,
    headers     bytea       NOT NULL,
    published   boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id)
)`
	CreateMessagesUnpublishedIndexSQL = `CREATE INDEX unpublished_idx ON messages (created_at) WHERE not published`
	CreateMessagesPublishedIndexSQL   = `CREATE INDEX published_idx ON messages (modified_at) WHERE published`

	CreateSagaInstancesTableSQL = `CREATE TABLE %s (
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
)`

	CreateSnapshotsTableSQL = `CREATE TABLE %s (
    entity_name      text        NOT NULL,
    entity_id        text        NOT NULL,
    snapshot_name    text        NOT NULL,
    snapshot_data    bytea       NOT NULL,
    snapshot_version int         NOT NULL,
    modified_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id)
)`

	loadEventsSQL = "SELECT event_name, event_data FROM %s WHERE entity_name = $1 AND entity_id = $2 AND event_version > $3 ORDER BY event_version ASC"
	writeEventSQL = "INSERT INTO %s (entity_name, entity_id, correlation_id, causation_id, event_version, event_name, event_data, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)"

	fetchMessagesSQL = "SELECT message_id, destination, payload, headers FROM %s WHERE not published ORDER BY created_at LIMIT %d"
	saveMessageSQL   = "INSERT INTO %s (message_id, destination, payload, headers) VALUES ($1, $2, $3, $4)"
	markMessagesSQL  = "UPDATE %s SET published = true, modified_at = CURRENT_TIMESTAMP WHERE message_id = ANY ($1)"
	purgeMessagesSQL = "DELETE FROM %s WHERE published AND modified_at < $1"

	findSagaInstanceSQL   = "SELECT saga_data_name, saga_data, current_step, end_state, compensating FROM %s WHERE saga_name = $1 AND saga_id = $2 LIMIT 1"
	saveSagaInstanceSQL   = "INSERT INTO %s (saga_name, saga_id, saga_data_name, saga_data, current_step, end_state, compensating, modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)"
	updateSagaInstanceSQL = "UPDATE %s SET saga_data = $1, current_step = $2, end_state = $3, compensating = $4, modified_at = CURRENT_TIMESTAMP WHERE saga_name = $5 AND saga_id = $6"

	loadSnapshotSQL = "SELECT snapshot_name, snapshot_data, snapshot_version FROM %s WHERE entity_name = $1 AND entity_id = $2 LIMIT 1"
	saveSnapshotSQL = `INSERT INTO %s (entity_name, entity_id, snapshot_name, snapshot_data, snapshot_version, modified_at) 
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) 
ON CONFLICT (entity_name, entity_id) DO
UPDATE SET snapshot_name = EXCLUDED.snapshot_name, snapshot_data = EXCLUDED.snapshot_data, snapshot_version = EXCLUDED.snapshot_version, modified_at = EXCLUDED.modified_at`

	pgxTxKey = contextKey(5432)
)

var ErrTxNotInContext = errors.New("pgx.Tx is not set for session")
var ErrInvalidTxValue = errors.New("tx value is not a pgx.Tx type")
