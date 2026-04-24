# 📁 Go RabbitMQ Patterns Masterclass

Este repositório reúne **projetos práticos em Go (Golang)** demonstrando os principais **padrões de mensageria com RabbitMQ**, com foco em **boas práticas**, organização de código e **cenários de carga/estresse** para validar comportamento e resiliência.

> Cada pasta de projeto possui seu próprio `README.md` com detalhes completos (comandos, filas/exchanges, rotas e exemplos).  
> **Use este README da raiz como visão geral e mapa do repositório.**

---

## ✅ O que você encontra aqui

Em sistemas distribuídos e microserviços, comunicação assíncrona é essencial. Aqui você aprende **como implementar** e **quando escolher** cada tipo de Exchange:

- **Direct Exchange**: roteamento **ponto a ponto** por chave exata (routing key).
- **Fanout Exchange**: **broadcast** (difusão) para todos os consumidores ligados.
- **Topic Exchange**: roteamento **seletivo** usando padrões com wildcard (`*` e `#`).

---

## 🧱 Estrutura (visão geral)

Os projetos seguem o **Standard Go Project Layout** (quando aplicável), normalmente com:

- `cmd/`: pontos de entrada (produtor, consumidor, stress test).
- `internal/`: configuração e infraestrutura (ex.: conexão/canal do RabbitMQ) encapsulada para reuso.
- `pkg/`: modelos/contratos compartilháveis (payloads, DTOs), quando fizer sentido.
- `docker-compose.yml`: infraestrutura para subir RabbitMQ local rapidamente.

> A estrutura exata pode variar por projeto — confira o `README.md` dentro de cada pasta.

---

## 🧪 Testes de carga/estresse (quando existir no projeto)

Alguns projetos incluem stress tests para simular alta concorrência com Go:

- disparos simultâneos com **goroutines**
- coordenação com **`sync.WaitGroup`**
- observação de **vazão (msg/s)** e comportamento pelo RabbitMQ Management

---

## 🛠 Pré-requisitos

- **Docker** + **Docker Compose**
- **Go 1.24+**
- (Opcional) Acesso ao **RabbitMQ Management** em `http://localhost:15672` (`guest/guest`)

---

## 🐇 Subir o RabbitMQ (ambiente local)

Na raiz do repositório (ou na pasta do projeto que contém o compose):

```bash
docker-compose up -d
```

Acesse o painel:

- `http://localhost:15672` (login `guest` / senha `guest`)

Para derrubar:

```bash
docker-compose down
```

---

## 🏁 Como rodar os exemplos

Cada projeto tem comandos próprios (nomes de binários, parâmetros e filas variam), então a forma correta é:

1. Entre na pasta do projeto desejado
2. Leia o `README.md` do projeto
3. Rode consumidores e produtores em terminais separados

### Fluxo típico (3 terminais)

- **Terminal 1**: consumidor(es)
- **Terminal 2**: produtor
- **Terminal 3**: (opcional) stress test / publicador em alta concorrência

> Nos READMEs de cada projeto, deixe explícito: exchange, filas, routing keys, bindings e exemplos de saída.

---

## 📊 Comparativo técnico (resumo)

- **Direct**: unicast → tarefas direcionadas, processamento por chave, logs específicos
- **Fanout**: broadcast → notificações, invalidação/atualização de cache, eventos “para todos”
- **Topic**: multicast seletivo → event bus com múltiplos tipos de evento e assinaturas por padrão

---

## 📚 Projetos neste repositório

- **[go-rabbitmq-direct](./go-rabbitmq-direct/README.md)**: Producer/Consumer usando **Direct Exchange** (`orders_direct`) com routing key `order.created` + stress test.
- **[go-rabbitmq-fanout](./go-rabbitmq-fanout/README.md)**: Pub/Sub usando **Fanout Exchange** (`orders_broadcast`) com múltiplos consumers (email/estoque).
- **[go-rabbitmq-topic](./go-rabbitmq-topic/README.md)**: **Topic Exchange** (`logs_topic`) com binds por padrão (`kern.*`, `#`) + stress test.

---