package broker

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	PublishOrderCreated(email, orderNumber string, totalAmount float64) error
}

type rabbitMQPublisher struct {
	conn *amqp.Connection
	url  string
}

func NewRabbitMQPublisher(url string) (EventPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &rabbitMQPublisher{conn: conn, url: url}, nil
}

func (p *rabbitMQPublisher) PublishOrderCreated(email, orderNumber string, totalAmount float64) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Pastikan exchange ada
	ch.ExchangeDeclare("ecommerce.events", "topic", true, false, false, false, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"email":   email,
		"subject": "Pesanan Dibuat: " + orderNumber,
		"body":    "Pesanan Anda berhasil dibuat. Silakan lakukan pembayaran.",
	}
	body, _ := json.Marshal(payload)

	return ch.PublishWithContext(ctx,
		"ecommerce.events", // exchange
		"order.created",    // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
