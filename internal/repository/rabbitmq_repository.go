package repository

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Publish(ctx context.Context, product interface{}) error
	AddReceiver(ctx context.Context, callback func(context.Context, amqp.Delivery) error)
	Close()
}

type queue struct {
	Channel  *amqp.Channel
	Name     string
	StopChan chan struct{}
}

func NewQueueDeclare(conn *amqp.Connection, queueName string) *queue {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if _, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		panic(err)
	}

	stopChan := make(chan struct{}, 1)

	return &queue{
		Channel:  ch,
		Name:     queueName,
		StopChan: stopChan,
	}
}

func (r *queue) Publish(ctx context.Context, product interface{}) error {
	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = r.Channel.Publish(
		"",     // default exchange
		r.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *queue) AddReceiver(ctx context.Context, callback func(context.Context, amqp.Delivery) error) {
	msgs, err := r.Channel.Consume(
		r.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				// Call the provided callback function to process the product
				errN := callback(ctx, msg)
				if errN != nil {
					msg.Nack(false, true) // Requeue the message
					continue
				}
				// Acknowledge message
				msg.Ack(true)
			case <-r.StopChan:
				log.Println("Stopping message consumption...")
				return
			}
		}
	}()
}

func (r *queue) Close() {
	close(r.StopChan)
	r.Channel.Close()
}
