package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ureis/go-rabbitmq-direct/internal/queue"
	"github.com/ureis/go-rabbitmq-direct/pkg/models"
)

func main() {
	// conexão com o RabbitMQ via Docker
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)

	}
	defer conn.Close()

	// Configurar a queue e a exchange
	ch, err := queue.SetupRabbitMQ(conn)
	if err != nil {
		log.Fatalf("Erro ao configurar o RabbitMQ: %v", err)
	}
	defer ch.Close()

	// Criar objeto de exemplo de pedido
	order := models.Order{
		ID:     "123",
		Amount: 250.00,
		Status: "PENDING",
	}

	body, _ := json.Marshal(order) // converte o objeto para JSON

	// Configurar o contexto e o cancelamento
	// O cancelamento é importante para evitar bloqueios de conexão
	// O timeout é importante para evitar que a conexão fique aberta para sempre
	// O context.Background() é o contexto raiz, e o WithTimeout é para cancelar a conexão após 5 segundos
	// O defer cancel() é importante para cancelar a conexão após a função terminar
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publicar a mensagem com a routing key correta "order.created"
	err = ch.PublishWithContext(ctx,
		"orders_direct", // exchange
		"order.created", // routing key (deve bater com o bind)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("Erro ao publicar: %v", err)
	}

	log.Printf("[x] Enviado: %s", body)

}
