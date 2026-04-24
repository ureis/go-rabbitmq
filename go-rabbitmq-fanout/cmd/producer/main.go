package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-fanout/pkg/models"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// No produtor, não precisamos criar fila, apenas a exchange
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(
		"orders_broadcast", // Nome da exchange
		"fanout",           // Tipo de exchange
		true,               // Durável
		false,              // Auto delete
		false,              // Interno
		false,              // Sem delay
		nil,                // Argumentos
	)

	order := models.Order{
		ID:     "FANOUT-999",
		Amount: 500.00,
		Status: "PAID",
	}
	body, _ := json.Marshal(order)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PUBLICANDO: Routing key vazia!
	err = ch.PublishWithContext(ctx,
		"orders_broadcast", // exchange
		"",                 // routing key (vazio para FANOUT)
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(" [x] Sent %s", body)

	ch.Close()
}
