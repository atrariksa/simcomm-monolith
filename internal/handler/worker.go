package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"log"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// type Queue interface {
// 	Publish(ctx context.Context, product interface{}) error
// 	Receive(ctx context.Context, queueName string, callback func(interface{}) error, stopChan <-chan struct{})
// }

// type queue struct {
// 	channel *amqp.Channel
// 	name    string
// }

// func NewQueueDeclare(conn *amqp.Connection, queueName string) *queue {
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if _, err := ch.QueueDeclare(
// 		queueName,
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	); err != nil {
// 		panic(err)
// 	}

// 	return &queue{
// 		channel: ch,
// 		name:    queueName,
// 	}
// }

// func (r *queue) Publish(product interface{}) error {
// 	body, err := json.Marshal(product)
// 	if err != nil {
// 		return err
// 	}

// 	err = r.channel.Publish(
// 		"",     // default exchange
// 		r.name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			DeliveryMode: amqp.Transient,
// 			ContentType:  "application/json",
// 			Body:         body,
// 		},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *queue) Receive(ctx context.Context, callback func(interface{}) error, stopChan <-chan struct{}) {
// 	msgs, err := r.channel.Consume(
// 		r.name,
// 		"",    // consumer
// 		false, // auto-ack
// 		false, // exclusive
// 		false, // no-local
// 		false, // no-wait
// 		nil,   // args
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register a consumer: %s", err)
// 	}

// 	go func() {
// 		for {
// 			select {
// 			case msg := <-msgs:
// 				// Call the provided callback function to process the product
// 				errN := callback(msg)
// 				if errN != nil {
// 					msg.Nack(false, true) // Requeue the message
// 					continue
// 				}
// 				// Acknowledge message
// 				msg.Ack(true)
// 			case <-stopChan:
// 				log.Println("Stopping message consumption...")
// 				return
// 			}
// 		}
// 	}()
// }

// func (r *queue) Close() {
// 	r.channel.Close()
// 	// r.connection.Close()
// }
