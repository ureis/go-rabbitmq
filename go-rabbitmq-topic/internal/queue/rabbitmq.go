package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupTopic(conn *amqp.Connection, queueName string, routingKey string) (*amqp.Channel, error) {
	ch, err := conn.Channel() // Cria o canal de comunicação com o RabbitMQ
	if err != nil {
		return nil, err
	}

	// 1. Exchange do tipo TOPIC
	err = ch.ExchangeDeclare( // Declara o exchange do tipo TOPIC
		"logs_topic", // nome do exchange
		"topic",      // tipo do exchange
		true,         // durable
		false,        // auto delete
		false,        // internal
		false,        // no wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	// 2. Declaração da Fila
	q, err := ch.QueueDeclare( // Declara a fila
		queueName, // nome da fila
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	// 3. Bind com Routing Key específico para a fila
	err = ch.QueueBind(
		q.Name,       // nome da fila
		routingKey,   // routing key
		"logs_topic", // nome do exchange
		false,        // no wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	return ch, err
}
