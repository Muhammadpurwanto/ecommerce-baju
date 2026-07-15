package broker

import (
	"encoding/json"
	"log"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn     *amqp.Connection
	emailSvc service.EmailService
}

func NewRabbitMQConsumer(url string, emailSvc service.EmailService) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{conn: conn, emailSvc: emailSvc}, nil
}

func (c *RabbitMQConsumer) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *RabbitMQConsumer) Listen() {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// Setup Exchange & Queues
	err = ch.ExchangeDeclare("ecommerce.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return
	}
	
	err = ch.ExchangeDeclare("user_events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare user_events exchange: %v", err)
		return
	}

	q, err := ch.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// Bind Queue ke routing keys
	ch.QueueBind(q.Name, "order.created", "ecommerce.events", false, nil)
	ch.QueueBind(q.Name, "payment.success", "ecommerce.events", false, nil)
	ch.QueueBind(q.Name, "user.registered", "user_events", false, nil)

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	log.Println("RabbitMQ Consumer started. Waiting for messages...")

	for d := range msgs {
		log.Printf("Received a message: %s", d.RoutingKey)

		// Ekstrak data (asumsi data JSON dengan field 'email' dan 'message')
		// Di implementasi nyata, buat DTO yang sesuai
		var payload map[string]interface{}
		if err := json.Unmarshal(d.Body, &payload); err == nil {
			email, _ := payload["email"].(string)

			var subject, body string
			if d.RoutingKey == "user.registered" {
				name, _ := payload["name"].(string)
				subject = "Selamat Datang di Ecommerce Baju!"
				body = "<h1>Halo " + name + "!</h1><p>Terima kasih telah mendaftar di layanan kami.</p>"
			} else {
				subject, _ = payload["subject"].(string)
				body, _ = payload["body"].(string)
			}

			if email != "" && subject != "" && body != "" {
				req := &dto.SendEmailRequest{
					To:      email,
					Subject: subject,
					Body:    body,
				}
				c.emailSvc.SendEmail(req)
			}
		}

		d.Ack(false) // Acknowledge pesan sukses diproses
	}
}
