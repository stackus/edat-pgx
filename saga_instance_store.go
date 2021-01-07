package edatpgx

import (
	"context"
	"fmt"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/saga"
)

type SagaInstanceStore struct {
	tableName string
	client    Client
	logger    log.Logger
}

var _ saga.InstanceStore = (*SagaInstanceStore)(nil)

func NewSagaInstanceStore(client Client, options ...SagaInstanceStoreOption) *SagaInstanceStore {
	s := &SagaInstanceStore{
		tableName: DefaultSagaInstanceTableName,
		client:    client,
		logger:    log.DefaultLogger,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *SagaInstanceStore) Find(ctx context.Context, sagaName, sagaID string) (*saga.Instance, error) {
	var dataName string
	var data []byte
	var currentStep int
	var endState, compensating bool

	row := s.client.QueryRow(ctx, fmt.Sprintf(findSagaInstanceSQL, s.tableName), sagaName, sagaID)
	err := row.Scan(&dataName, &data, &currentStep, &endState, &compensating)
	if err != nil {
		return nil, err
	}

	sagaData, err := core.DeserializeSagaData(dataName, data)
	if err != nil {
		return nil, err
	}

	return saga.NewSagaInstance(sagaName, sagaID, sagaData, currentStep, endState, compensating), nil
}

func (s *SagaInstanceStore) Save(ctx context.Context, sagaInstance *saga.Instance) error {
	data, err := core.SerializeSagaData(sagaInstance.SagaData())
	if err != nil {
		return err
	}
	_, err = s.client.Exec(ctx, fmt.Sprintf(saveSagaInstanceSQL, s.tableName), sagaInstance.SagaName(), sagaInstance.SagaID(), sagaInstance.SagaData().SagaDataName(), data, sagaInstance.CurrentStep(), sagaInstance.EndState(), sagaInstance.Compensating())
	return err
}

func (s *SagaInstanceStore) Update(ctx context.Context, sagaInstance *saga.Instance) error {
	data, err := core.SerializeSagaData(sagaInstance.SagaData())
	if err != nil {
		return err
	}
	_, err = s.client.Exec(ctx, fmt.Sprintf(updateSagaInstanceSQL, s.tableName), data, sagaInstance.CurrentStep(), sagaInstance.EndState(), sagaInstance.Compensating(), sagaInstance.SagaName(), sagaInstance.SagaID())
	return err
}
