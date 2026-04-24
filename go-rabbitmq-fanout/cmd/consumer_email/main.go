package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-fanout/internal/queue"
	"github.com/ureis/go-rabbitmq-fanout/pkg/models"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := queue.SetupFanout(conn, "queue_email")
	if err != nil {
		log.Fatalf("Erro ao configurar o RabbitMQ: %v", err)
	}
	defer ch.Close()

	// Inicar o consumo da fila
	msgs, err := ch.Consume(
		"queue_email", // Nome da queue(fila)
		"",            // Consumer tag
		true,          // auto-ack (confirmação automática de recebimento)
		false,         // exclusive (exclusivo para este consumidor)
		false,         // no-local (não receber mensagens enviadas pelo mesmo consumidor)
		false,         // no-wait (não aguarda resposta)
		nil,           // arguments (não usados)
	)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var order models.Order
			if err := json.Unmarshal(d.Body, &order); err != nil {
				log.Printf("Erro ao decodificar a mensagem: %v", err)
				continue
			}
			log.Printf(" [v] Recebido com sucesso: %v", order)
		}
	}()

	log.Printf(" [*] Aguardando mensagens. Para sair, pressione CTRL+C")
	<-forever
}
