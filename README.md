# edat-pgx - Postgres stores for edat

## Installation

    go get -u github.com/stackus/edat-pgx

## Usage Example

    import "github.com/stack/edat-pgx"

    conn, _ := pgxpool.Connect(ctx, "your-connection-string")

    // Create a store for aggregate events using the pool connection
    eventStore := edatpgx.NewEventStore(conn)

    // Create a store using a session client that uses a pgx.Tx from context
    client := edatpgx.NewSessionClient()
    eventStore := edatpgx.NewEventStore(client)


## Prerequisites

Go 1.16

## Features

Stores accept `*pgx.Conn`, `*pgxpool.Pool`, `pgx.Tx` and `edatpgx.Client` for clients. Middleware will accept `*pgxpool.Pool` only.

- Session Client `NewSessionClient()`
- Aggregate Event Store `NewEventStore(client, ...options)`
- Outbox Message Store and Producer `NewMessageStore(client, ....options)`
- Saga Instance Store `NewSagaInstanceStore(client, ...options)`
- Aggregate Snapshot Store `NewSnapshotStore(client, ...options)`
- Message Receiver Session Middleware `ReceiverSessionMiddleware(*pgxpool.Pool, log.Logger)`
- Web Request Session Middleware `WebSessionMiddleware(*pgxpool.Pool, log.Logger)`
- Grpc Request Session (Unary) Interceptor `RpcSessionUnaryMiddleware(*pgxpool.Pool, log.Logger)`

## TODOs

- Documentation
- Tests, tests, and more tests

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT
