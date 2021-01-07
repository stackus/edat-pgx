package edatpgx

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client covers a subset of what both pgx.Conn or pgxpool.Pool provide
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// This is unnecessary right now
//
// type client struct {
// 	conn *pgxpool.Pool
// }
//
// var _ Client = (*client)(nil)
//
// func NewClient(conn *pgxpool.Pool) Client {
// 	return &client{conn}
// }
//
// func (c client) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
// 	return c.conn.Exec(ctx, sql, arguments...)
// }
//
// func (c client) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
// 	return c.conn.Query(ctx, sql, args...)
// }
//
// func (c client) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
// 	return c.conn.QueryRow(ctx, sql, args...)
// }
//
// func (c client) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
// 	return c.conn.QueryFunc(ctx, sql, args, scans, f)
// }
//
// func (c client) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
// 	return c.conn.SendBatch(ctx, b)
// }
//
// func (c client) Begin(ctx context.Context) (pgx.Tx, error) {
// 	return c.conn.Begin(ctx)
// }
//
// func (c client) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
// 	return c.conn.BeginTx(ctx, txOptions)
// }

type sessionClient struct{}

var _ Client = (*sessionClient)(nil)

// NewSessionClient returns a pgx.Conn or pgxpool.Pool compatible client that uses an active transaction from context
func NewSessionClient() Client {
	return sessionClient{}
}

func (c sessionClient) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	tx, err := c.tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Exec(ctx, sql, arguments...)
}

func (c sessionClient) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	tx, err := c.tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Query(ctx, sql, args...)
}

func (c sessionClient) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	tx, err := c.tx(ctx)
	if err != nil {
		return rowError{err}
	}

	return tx.QueryRow(ctx, sql, args...)
}

func (c sessionClient) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	tx, err := c.tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.QueryFunc(ctx, sql, args, scans, f)
}

func (c sessionClient) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	tx, err := c.tx(ctx)
	if err != nil {
		return batchError{err}
	}

	return tx.SendBatch(ctx, b)
}

func (c sessionClient) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := c.tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Begin(ctx)
}

func (c sessionClient) BeginTx(ctx context.Context, _ pgx.TxOptions) (pgx.Tx, error) {
	tx, err := c.tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Begin(ctx)
}

func (s sessionClient) tx(ctx context.Context) (pgx.Tx, error) {
	value := ctx.Value(pgxTxKey)
	if value == nil {
		return nil, ErrTxNotInContext
	}

	tx, ok := value.(pgx.Tx)
	if !ok {
		return nil, ErrInvalidTxValue
	}

	return tx, nil
}

type rowError struct {
	err error
}

func (e rowError) Scan(...interface{}) error {
	return e.err
}

type batchError struct {
	err error
}

var _ pgx.BatchResults = (*batchError)(nil)

func (e batchError) Exec() (pgconn.CommandTag, error) {
	return nil, e.err
}

func (e batchError) Query() (pgx.Rows, error) {
	return nil, e.err
}

func (e batchError) QueryRow() pgx.Row {
	return rowError{e.err}
}

func (e batchError) Close() error {
	return e.err
}
