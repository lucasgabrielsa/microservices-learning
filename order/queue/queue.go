package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func Connect() *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")

	conn, error := amqp.Dial(dsn)
	if error != nil {
		panic(error.Error())
	}

	channel, error := conn.Channel()
	if error != nil {
		panic(error.Error())
	}

	return channel
}

func Notify(payload []byte, exchange string, routingKey string, channel *amqp.Channel) {

	error := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: []byte(payload),
		},
	)

	if error != nil {
		panic(error.Error())
	}

	fmt.Println("Message sent")
}

func StartConsuming(channel *amqp.Channel, in chan []byte) {
	q, error := channel.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
		)

	if error != nil {
		panic(error.Error())
	}

	msgs, error := channel.Consume(
		q.Name,
		"checkout",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for m:= range msgs {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}