CREATE TABLE events
(
    entity_name    text        NOT NULL,
    entity_id      text        NOT NULL,
    correlation_id text        NOT NULL,
    causation_id   text        NOT NULL,
    event_version  int         NOT NULL,
    event_name     text        NOT NULL,
    event_data     bytea       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id, event_version)
);

CREATE TABLE messages
(
    message_id  text        NOT NULL,
    destination text        NOT NULL,
    payload     bytea       NOT NULL,
    headers     bytea       NOT NULL,
    published   boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id)
);

CREATE INDEX unpublished_idx ON messages (created_at) WHERE NOT published;
CREATE INDEX published_idx ON messages (modified_at) WHERE published;

CREATE TABLE saga_instances
(
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
);

CREATE TABLE snapshots
(
    entity_name      text        NOT NULL,
    entity_id        text        NOT NULL,
    snapshot_name    text        NOT NULL,
    snapshot_data    bytea       NOT NULL,
    snapshot_version int         NOT NULL,
    modified_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entity_name, entity_id)
);
