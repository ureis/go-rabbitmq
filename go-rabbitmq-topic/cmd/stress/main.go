package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const totalMessages = 1000
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	// Garantir que a exchange existe
	_ = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)

	var wg sync.WaitGroup
	start := time.Now()

	// Lista de chaves para testar diferentes comportamentos
	keys := []string{"kern.info", "kern.error", "auth.info", "auth.error", "app.debug"}

	fmt.Printf("🚀 Iniciando stress test TOPIC com %d mensagens...\n", totalMessages)

	for i := 0; i < totalMessages; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Escolhe uma chave da lista baseado no índice
			routingKey := keys[id%len(keys)]

			body := fmt.Sprintf("ID:%d | Log: Evento ocorrido em %s", id, routingKey)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err := ch.PublishWithContext(ctx,
				"logs_topic",
				routingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})

			if err != nil {
				log.Printf("Erro: %v", err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n✅ Stress finalizado em %v\n", time.Since(start))
}
