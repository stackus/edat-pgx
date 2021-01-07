package edatpgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"
)

type MessageStore struct {
	tableName string
	client    Client
	logger    log.Logger
}

var _ msg.Producer = (*MessageStore)(nil)
var _ outbox.MessageStore = (*MessageStore)(nil)

func NewMessageStore(client Client, options ...MessageStoreOption) *MessageStore {
	s := &MessageStore{
		tableName: DefaultMessageTableName,
		client:    client,
		logger:    log.DefaultLogger,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *MessageStore) Send(ctx context.Context, channel string, message msg.Message) error {
	headers, err := json.Marshal(message.Headers())
	if err != nil {
		return err
	}

	return s.Save(ctx, outbox.Message{
		MessageID:   message.ID(),
		Destination: channel,
		Payload:     message.Payload(),
		Headers:     headers,
	})
}

func (s *MessageStore) Close(ctx context.Context) error {
	return nil
}

func (s *MessageStore) Fetch(ctx context.Context, limit int) ([]outbox.Message, error) {
	messages := make([]outbox.Message, 0)

	rows, err := s.client.Query(ctx, fmt.Sprintf(fetchMessagesSQL, s.tableName, limit))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return messages, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message outbox.Message

		err := rows.Scan(&message.MessageID, &message.Destination, &message.Payload, &message.Headers)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	return messages, nil
}

func (s *MessageStore) Save(ctx context.Context, message outbox.Message) error {
	_, err := s.client.Exec(ctx, fmt.Sprintf(saveMessageSQL, s.tableName),
		message.MessageID, message.Destination, message.Payload, message.Headers,
	)
	return err
}

func (s *MessageStore) MarkPublished(ctx context.Context, messageIDs []string) error {
	ids := &pgtype.TextArray{}
	err := ids.Set(messageIDs)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, fmt.Sprintf(markMessagesSQL, s.tableName), ids)
	return err
}

func (s *MessageStore) PurgePublished(ctx context.Context, olderThan time.Duration) error {
	var when time.Time

	// cannot decide if I want positive or negative durations
	if olderThan < 0 {
		when = time.Now().Add(olderThan)
	} else {
		when = time.Now().Add(-1 * olderThan)
	}

	_, err := s.client.Exec(ctx, fmt.Sprintf(purgeMessagesSQL, s.tableName), when)
	return err
}
