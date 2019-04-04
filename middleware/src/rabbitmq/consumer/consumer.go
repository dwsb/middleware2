package consumer

import (
	"errors"
	"middleware2/middleware/src/utils"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel    *amqp.Channel
	queue      amqp.Queue
	queueReply amqp.Queue
}

func (c *Consumer) Connect(name string) error {
	ch, qu, qur, err := utils.ConnectRabbitMQ(name)
	if err != nil {
		return err
	}

	c.channel = ch
	c.queue = qu
	c.queueReply = qur

	return nil
}

func (c *Consumer) Consume(contentType string) ([]byte, error) {
	if c.channel == nil {
		return nil, errors.New("Canal não está aberto")
	}

	var delivery <-chan amqp.Delivery
	delivery, err := utils.ConsumeQueue(c.queueReply.Name, c.channel)
	if err != nil {
		return nil, err
	}

	reply := <-delivery
	return reply.Body, nil
}

func (c *Consumer) Reply(message []byte, contentType string) error {
	return utils.PublishQueue(c.queue.Name, contentType, message, c.channel)
}
