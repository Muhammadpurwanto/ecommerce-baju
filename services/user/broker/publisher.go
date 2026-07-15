package broker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	PublishUserRegistered(uuid string, email string, name string) error
}

type rabbitMQPublisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitMQPublisher(amqpURL string) (EventPublisher, error) {
	var conn *amqp091.Connection
	var err error

	for i := 0; i < 5; i++ { 
		conn, err = amqp091.Dial(amqpURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (User Publisher), retrying in 5 seconds... (%v)", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare( 
		"user_events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (p *rabbitMQPublisher) PublishUserRegistered(uuid string, email string, name string) error {
	payload := map[string]string{
		"uuid":  uuid,
		"email": email,
		"name":  name,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(
		context.Background(),	
		"user_events",     // exchange
		"user.registered", // routing key
		false,             // mandatory
		false,             // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish UserRegistered event: %v", err)
		return err
	}

	log.Printf("Published UserRegistered event for user %s", uuid)
	return nil
}
