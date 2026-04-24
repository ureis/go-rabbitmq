package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-topic/internal/queue"
)

func main() {
	// Conexão com o RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Configuração do Canal e bind
	// Nome da fila: "queue_kernel"
	// Routing Key: "kern.* (O asterisco substitui uma palavra qualquer)"
	ch, err := queue.SetupTopic(conn, "queue_kernel", "kern.*")
	if err != nil {
		log.Fatal("Erro ao configurar Topic: ", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"queue_kernel", // nome da fila
		"",             // consumer
		true,           // auto ack
		false,          // exclusive
		false,          // no local
		false,          // no wait
		nil,            // arguments
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [KERNEL MONITOR] Recebido: %s | Key: %s", d.Body, d.RoutingKey)
		}
	}()

	log.Printf(" [*] Aguardando logs de KERNEL (kern.*). Para sair pressione CTRL+C")

	<-forever
}
