package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupRabbitMQ configura a conexão com o RabbitMQ
func SetupRabbitMQ(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 1. Declarar a Exchange do tipo Direct
	// Por que: Garante que a exchange exista antes de enviar/receber.
	err = ch.ExchangeDeclare(
		"orders_direct", // Nome da exchange
		"direct",        // Tipo de exchange
		true,            // Durable (sobrevive ao restart do broker)
		false,           // Auto delete (apaga se não tiver consumidores)
		false,           // Internal (não é usado para roteamento interno)
		false,           // No-wait (não aguarda resposta)
		nil,             // Arguments (não usados)
	)

	if err != nil {
		return nil, err
	}

	// 2. Declarar a Queue(fila)
	q, err := ch.QueueDeclare(
		"orders_queue", // Nome da queue
		true,           // Durable (sobrevive ao restart do broker)
		false,          // Auto delete (apaga se não tiver consumidores)
		false,          // Exclusivo (apenas o primeiro consumidor)
		false,          // No-wait (não aguarda resposta)
		nil,            // Arguments (não usados)
	)

	if err != nil {
		return nil, err
	}

	// 3. Bind: Ligar o Queue(fila) à Exchange usando uma Routing Key
	// Por que: No modo Direct, a mensagem só cai na fila se a routing key for igual à key da exchange.
	err = ch.QueueBind(
		q.Name,          // Nome da queue(fila)
		"order.created", // Routing Key
		"orders_direct", // Nome da exchange
		false,           // No-wait (não aguarda resposta)
		nil,             // Arguments (não usados)
	)

	return ch, err
}
