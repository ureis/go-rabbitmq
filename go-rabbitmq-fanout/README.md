# go-rabbitmq-fanout

Exemplo em Go de **pub/sub com RabbitMQ usando Exchange do tipo FANOUT**.  
Um producer publica pedidos na exchange `orders_broadcast` e múltiplos consumers recebem a mesma mensagem em filas diferentes.

## Arquitetura

- **Exchange (fanout)**: `orders_broadcast`
- **Filas (uma por consumer)**:
  - `queue_email` (consumer `cmd/consumer_email`)
  - `queue_estoque` (consumer `cmd/consumer_estoque`)
- **Mensagem**: `pkg/models.Order`

No FANOUT a routing key é ignorada (por isso o publish usa `""`).

## Requisitos

- Go **1.24+**
- Docker + Docker Compose

## Subindo o RabbitMQ

Na raiz do projeto:

```bash
docker compose up -d
```

- AMQP: `localhost:5672`
- Management UI: `http://localhost:15672` (user `guest`, pass `guest`)

## Rodando os consumers

Abra **dois terminais** (um para cada consumer):

```bash
go run cmd/consumer_email/main.go
```

```bash
go run cmd/consumer_estoque/main.go
```

Eles criam/bindam suas filas na exchange `orders_broadcast` via `internal/queue.SetupFanout`.

## Rodando o producer

Em outro terminal:

```bash
go run cmd/producer/main.go
```

Você deve ver a mensagem chegando **nos dois consumers**.

## Troubleshooting

- **`channel/connection is not open`**
  - Em geral isso acontece quando o channel foi fechado antes de `ExchangeDeclare`/`Publish` ou quando o RabbitMQ não está no ar.
  - Verifique se o RabbitMQ está rodando com `docker compose ps` e se a URL AMQP está correta: `amqp://guest:guest@localhost:5672/`.

## Estrutura de pastas

- `cmd/producer`: publica uma `Order` na exchange fanout
- `cmd/consumer_email`: consome da fila `queue_email`
- `cmd/consumer_estoque`: consome da fila `queue_estoque`
- `internal/queue`: setup da exchange/queue/bind (`SetupFanout`)
- `pkg/models`: modelos compartilhados (ex.: `Order`)

