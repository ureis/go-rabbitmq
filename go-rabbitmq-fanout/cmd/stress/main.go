package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-fanout/pkg/models"
)

func main() {
	const totalMessages = 1000 // Quantidade de mensagens

	// Conexão com o RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Falha na conexão: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir canal: %v", err)
	}
	defer ch.Close()

	// Garantimos que a exchange existe
	err = ch.ExchangeDeclare(
		"orders_broadcast",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🔥 Iniciando Stress Test FANOUT: %d mensagens...\n", totalMessages)

	for i := 1; i <= totalMessages; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			order := models.Order{
				ID:     fmt.Sprintf("STRESS-FANOUT-%d", id),
				Amount: float64(id) * 2.0,
				Status: "BULK_BROADCAST",
			}

			body, _ := json.Marshal(order)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// No Fanout, a Routing Key (segundo parâmetro) DEVE ser vazia ""
			err := ch.PublishWithContext(ctx,
				"orders_broadcast",
				"",
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent, // Persistência em disco
					ContentType:  "application/json",
					Body:         body,
				})

			if err != nil {
				log.Printf("❌ Erro no envio %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("\n✅ Stress Test concluído!")
	fmt.Printf("\n⏱  Tempo de envio: %v", duration)
	fmt.Printf("\n📢 Mensagens enviadas: %d", totalMessages)
	fmt.Printf("\n📈 Vazão de saída: %.2f msg/s\n", float64(totalMessages)/duration.Seconds())
	fmt.Println("Observe no Management a duplicação nas filas de estoque e email.")
}
