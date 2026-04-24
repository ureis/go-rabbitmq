package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ureis/go-rabbitmq-topic/internal/queue"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Configuração do Canal e Bind
	// Nome da fila: "queue_all_logs"
	// Routing Key: "#" (Recebe TUDO de qualquer nível)
	ch, err := queue.SetupTopic(conn, "queue_all_logs", "#")
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"queue_all_logs", // nome da fila
		"",               // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              // arguments
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [GLOBAL MONITOR] Recebido: %s | Key: %s", d.Body, d.RoutingKey)
		}
	}()

	log.Printf(" [*] Aguardando TODOS os logs (#). Para sair pressione CTRL+C")

	<-forever

}
