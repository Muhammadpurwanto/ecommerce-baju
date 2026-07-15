package broker

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	PublishPaymentSuccess(email, orderID string) error
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

func (p *rabbitMQPublisher) PublishPaymentSuccess(email, orderID string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	ch.ExchangeDeclare("ecommerce.events", "topic", true, false, false, false, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"email":    email,
		"subject":  "Pembayaran Berhasil untuk Order " + orderID,
		"body":     "Pembayaran Anda telah kami terima. Pesanan sedang diproses.",
		"order_id": orderID,
	}
	body, _ := json.Marshal(payload)

	return ch.PublishWithContext(ctx,
		"ecommerce.events",
		"payment.success", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
