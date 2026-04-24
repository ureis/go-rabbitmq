package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-direct/internal/queue"
	"github.com/ureis/go-rabbitmq-direct/pkg/models"
)

func main() {
	const totalMessages = 10000 // Quantidade de mensagens a serem enviadas

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := queue.SetupRabbitMQ(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Iniciando stress test: %d mensagens...\n", totalMessages)

	for i := 1; i <= totalMessages; i++ {
		wg.Add(1)

		// Dispara uma Goroutine por mensagem
		go func(id int) {
			defer wg.Done()

			order := models.Order{
				ID:     fmt.Sprintf("STRESS-%d", id),
				Amount: float64(id) * 1.5,
				Status: "STRESS_TEST",
			}

			body, _ := json.Marshal(order)

			// Contexto curto para casa publicação
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err := ch.PublishWithContext(ctx,
				"orders_direct",
				"order.created",
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent, // Garante que a mensagem seja persistente
					ContentType:  "application/json",
					Body:         body,
				})

			if err != nil {
				log.Printf("[x] Erro no ID %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait() // Aguarda todas as Goroutines terminarem
	duration := time.Since(start)

	fmt.Printf("\n✅ Teste finalizado!")
	fmt.Printf("\n⏱  Tempo total: %v", duration)
	fmt.Printf("\n📈 Média: %.2f msg/s\n", float64(totalMessages)/duration.Seconds())
}
