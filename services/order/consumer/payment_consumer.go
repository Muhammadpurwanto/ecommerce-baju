package consumer

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn     *amqp.Connection
	orderSvc service.OrderService
}

func NewRabbitMQConsumer(url string, orderSvc service.OrderService) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{conn: conn, orderSvc: orderSvc}, nil
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

	// Setup Exchange & Queue
	err = ch.ExchangeDeclare("ecommerce.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return
	}

	q, err := ch.QueueDeclare("order_payment_queue", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// Bind Queue ke routing key payment.success
	err = ch.QueueBind(q.Name, "payment.success", "ecommerce.events", false, nil)
	if err != nil {
		log.Printf("Failed to bind queue: %v", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	log.Println("RabbitMQ Consumer (Order) started. Waiting for payment success messages...")

	for d := range msgs {
		log.Printf("Received payment event: %s", d.RoutingKey)

		var payload map[string]interface{}
		if err := json.Unmarshal(d.Body, &payload); err == nil {
			orderIDVal, exists := payload["order_id"]
			if exists {
				var orderID uint64
				switch v := orderIDVal.(type) {
				case string:
					orderID, _ = strconv.ParseUint(v, 10, 64)
				case float64:
					orderID = uint64(v)
				}

				if orderID > 0 {
					log.Printf("Updating Order ID %d to paid", orderID)
					_, err := c.orderSvc.UpdateOrderStatus(uint(orderID), "paid")
					if err != nil {
						log.Printf("Failed to update order status for order %d: %v", orderID, err)
					} else {
						log.Printf("Successfully updated Order ID %d status to paid", orderID)
					}
				}
			}
		}

		d.Ack(false)
	}
}
