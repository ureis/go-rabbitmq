package main

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()

	_ = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	// Vamos enviar mensagens para diferentes "topicos"
	keys := []string{
		"kern.info",
		"kern.error",
		"auth.info",
		"auth.error",
	}

	for _, key := range keys {
		message := fmt.Sprintf("Log message with key: %s", key)

		_ = ch.PublishWithContext(
			context.Background(),
			"logs_topic",
			key,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)

		log.Printf(" [x] Enviado: %s", key)
	}
}
