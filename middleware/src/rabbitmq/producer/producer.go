package producer

import (
	"errors"
	"fmt"
	"middleware2/middleware/src/utils"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel    *amqp.Channel
	queue      amqp.Queue
	queueReply amqp.Queue
}

func (p *Producer) Connect(name string) error {
	ch, qu, qur, err := utils.ConnectRabbitMQ(name)
	if err != nil {
		return err
	}

	p.channel = ch
	p.queue = qu
	p.queueReply = qur

	return nil
}

func (p *Producer) Close() error {
	return p.channel.Close()
}

func (p *Producer) ProduceAndWaitReply(message []byte, contentType string) ([]byte, error) {
	if p.channel == nil {
		return nil, errors.New("Connect must be called before")
	}

	err := utils.PublishQueue(p.queue.Name, contentType, message, p.channel)

	if err != nil {
		return nil, err
	}

	var delivery <-chan amqp.Delivery
	delivery, err = utils.ConsumeQueue(p.queueReply.Name, p.channel)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	reply := <-delivery
	return reply.Body, nil
}
