package queue

import amqp "github.com/rabbitmq/amqp091-go"

func SetupFanout(conn *amqp.Connection, queueName string) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 1. Declare a Exchange como FANOUT
	err = ch.ExchangeDeclare(
		"orders_broadcast", // Nome da exchange
		"fanout",           // Tipo de exchange
		true,               // Durável
		false,              // Auto delete
		false,              // Interno
		false,              // Sem delay
		nil,                // Argumentos
	)
	if err != nil {
		return nil, err
	}

	// 2. Declarar a Fila especifica para este consumidor
	q, err := ch.QueueDeclare(
		queueName, // Nome da fila (cada consumer tem sua própria fila)
		true,      // Durável
		false,     // Excluível
		false,     // Exclusiva
		false,     // Sem delay
		nil,       // Argumentos
	)
	if err != nil {
		return nil, err
	}

	// 3. Bind SEM Routing Key
	// No Fanout, a routing key é ignorada ("")
	err = ch.QueueBind(
		q.Name,             // Nome da fila
		"",                 // Routing key (ignorado para FANOUT)
		"orders_broadcast", // Nome da exchange
		false,              // Sem delay
		nil,                // Argumentos
	)
	if err != nil {
		return nil, err
	}
	return ch, err
}
