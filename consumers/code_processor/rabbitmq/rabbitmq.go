package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQReceiver struct {
    conn *amqp.Connection
    ch *amqp.Channel
    queue *amqp.Queue
}

func NewRabbitMQReceiver(url string) (*RabbitMQReceiver, error) {
    conn, err := amqp.Dial(url)

    if err != nil {
        return nil, fmt.Errorf("connecting rabbitmq: %s", err.Error())
    }

    ch, err := conn.Channel()

    if err != nil {
        return nil, fmt.Errorf("opening rabbitmq channel: %s", err.Error())
    }

    queue, err := ch.QueueDeclare(
        "code_transfer",
        true,
        false,
        false,
        false,
        nil,
    )

    if err != nil {
        return nil, fmt.Errorf("declaring rabbitmq queue: %s", err.Error())
    }

    return &RabbitMQReceiver{
        conn: conn,
        ch: ch,
        queue: &queue,
    }, nil
}

func (s *RabbitMQReceiver) StartReceive() error {
    msgs, err := s.ch.Consume(
        s.queue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )

    if err != nil {
        return fmt.Errorf("rabbitmq, starting consuming: %s", err.Error())
    }

    go func() {
        for m := range msgs {
            log.Printf("message: %s", m.Body)
        }
    }()

    return nil
}

func (s *RabbitMQReceiver) Close() {
    s.ch.Close()
    s.conn.Close()
}