# go-rabbitmq-direct

Exemplo simples de **Producer/Consumer em Go** usando **RabbitMQ** com **Direct Exchange**.

## Arquitetura (Direct)

- **Exchange**: `orders_direct` (tipo `direct`)
- **Queue**: `orders_queue`
- **Binding / Routing key**: `order.created`

No modo **direct**, a mensagem só chega na fila quando a **routing key** do publish é **igual** à routing key do bind.

## Requisitos

- **Go** `1.24+` (ver `go.mod`)
- **Docker** (para subir o RabbitMQ)

## Subindo o RabbitMQ

Na raiz do projeto:

```bash
docker compose up -d
```

RabbitMQ Management:

- URL: `http://localhost:15672`
- Usuário/Senha: `guest` / `guest`

AMQP:

- Host: `localhost`
- Porta: `5672`
- URL: `amqp://guest:guest@localhost:5672/`

## Rodando o consumer

Em um terminal:

```bash
go run ./cmd/consumer
```

Ele fica aguardando mensagens e imprime o payload decodificado do JSON.

## Enviando uma mensagem (producer)

Em outro terminal:

```bash
go run ./cmd/producer
```

O producer publica um pedido JSON na exchange `orders_direct` com a routing key `order.created`.

Payload de exemplo:

```json
{"id":"123","amount":250,"status":"PENDING"}
```

## Stress test

Para enviar muitas mensagens em paralelo:

```bash
go run ./cmd/stress
```

Por padrão envia **10.000** mensagens (ajuste em `cmd/stress/main.go`).

## Estrutura do projeto

- `cmd/producer`: publica mensagem na exchange
- `cmd/consumer`: consome mensagens da fila
- `cmd/stress`: dispara muitas publicações em paralelo
- `internal/queue`: declara exchange/queue e faz o bind
- `pkg/models`: modelos compartilhados (ex.: `Order`)

## Notas

- A declaração de exchange/queue/bind acontece em `internal/queue/SetupRabbitMQ`.
- O consumer está com **auto-ack** ligado (`Consume(..., autoAck=true, ...)`). Para cenários reais, normalmente você desliga e dá ack manual após processar com sucesso.

