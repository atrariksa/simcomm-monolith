package util

import (
	"fmt"
	"log"
	"simcomm-monolith/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetRabbitMQConnection(cfg config.RabbitMQConfig) (connection *amqp.Connection) {
	host := cfg.Host
	user := cfg.User
	password := cfg.Password
	dsn := fmt.Sprintf("amqp://%v:%v@%v", user, password, host)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
