# go-rabbitmq-topic

Exemplo em Go de **RabbitMQ Topic Exchange** usando `github.com/rabbitmq/amqp091-go`.

O projeto publica mensagens no exchange `logs_topic` (tipo `topic`) com diferentes *routing keys* e possui consumidores que fazem *bind* com padrões (ex.: `kern.*`, `#`).

## Requisitos

- **Go**: 1.24.1+
- **RabbitMQ** rodando localmente em `amqp://guest:guest@localhost:5672/`

## Subindo o RabbitMQ (Docker)

```bash
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

- UI de management: `http://localhost:15672` (user/pass: `guest`/`guest`)

## Como executar

Na raiz do projeto:

### Produtor

Publica mensagens no exchange `logs_topic` para algumas chaves.

```bash
go run ./cmd/producer
```

**Routing keys enviadas pelo produtor (conforme código):**
- `kern.info`
- `kern.error`
- `auth.info`
- `auth.error`

### Consumidor (Kernel)

Cria a fila `queue_kernel` e faz bind com `kern.*`.

```bash
go run ./cmd/consumer_kernel
```

### Consumidor (Todos os logs)

Cria a fila `queue_all_logs` e faz bind com `#` (recebe tudo).

```bash
go run ./cmd/consumer_todos
```

### Stress test

Dispara 1000 publicações concorrentes no exchange `logs_topic`, alternando as chaves:
`kern.info`, `kern.error`, `auth.info`, `auth.error`, `app.debug`.

```bash
go run ./cmd/stress
```

## Notas importantes

- O exchange `logs_topic` é declarado como `durable=true`.
- Os consumidores usam `autoAck=true`.
- Os nomes e routing keys foram padronizados para `kern.*` e `queue_kernel` para evitar confusão.

## Estrutura

- `cmd/producer`: publica mensagens com routing keys
- `cmd/consumer_kernel`: consome `kern.*`
- `cmd/consumer_todos`: consome `#` (tudo)
- `cmd/stress`: publica em alta concorrência
- `internal/queue`: helper `SetupTopic` (exchange/queue/bind)

