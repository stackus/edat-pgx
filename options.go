package edatpgx

import (
	"github.com/stackus/edat/log"
)

type LoggerOption struct {
	log.Logger
}

func WithLogger(s log.Logger) LoggerOption {
	return LoggerOption{s}
}

func (o LoggerOption) configureEventStore(s *EventStore) {
	s.logger = o
}

func (o LoggerOption) configureMessageStore(s *MessageStore) {
	s.logger = o
}

func (o LoggerOption) configureSagaInstanceStore(s *SagaInstanceStore) {
	s.logger = o
}

func (o LoggerOption) configureSnapshotStore(s *SnapshotStore) {
	s.logger = o
}

func (o LoggerOption) configureSessionMiddlewareStore(c *SessionMiddlewareConfig) {
	c.logger = o
}

type TableNameOption string

func WithTableName(tn string) TableNameOption {
	return TableNameOption(tn)
}

func (o TableNameOption) configureEventStore(s *EventStore) {
	s.tableName = string(o)
}

func (o TableNameOption) configureMessageStore(s *MessageStore) {
	s.tableName = string(o)
}

func (o TableNameOption) configureSagaInstanceStore(s *SagaInstanceStore) {
	s.tableName = string(o)
}

func (o TableNameOption) configureSnapshotStore(s *SnapshotStore) {
	s.tableName = string(o)
}
