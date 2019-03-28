package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/streadway/amqp"
)

func Encode(bytes []byte) []byte {
	return []byte(string(bytes) + "\n")
}

func EncodeString(value string) []byte {
	return []byte(value + "\n")
}

func OpenConnection(protocol, port string) (net.Conn, error) {
	return net.Dial(protocol, fmt.Sprintf(":%s", port))
}

func FormatString(s string) string {
	s = strings.Replace(s, "\n", "", 1)
	s = strings.Replace(s, "\r", "", 1)

	return s
}

func ConnectRabbitMQ() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:1234")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	return ch, err
}

func DeclareQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare("teste", false, false, false, false, nil)
}

func ConsumeQueue(name string, ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	return ch.Consume(name, "", true, false, false, false, nil)
}

func PublishQueue(name, contentType string, body []byte, ch *amqp.Channel) error {
	return ch.Publish("",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
}
