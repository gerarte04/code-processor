package rabbitmq

import (
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQReceiver struct {
    conn *amqp.Connection
    ch *amqp.Channel
    queue *amqp.Queue

    msgHandler usecases.MessageHandler
}

func NewRabbitMQReceiver(url string, msgHandler usecases.MessageHandler) (*RabbitMQReceiver, error) {
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

        msgHandler: msgHandler,
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

    var forever chan struct{}

    go func() {
        for m := range msgs {
            log.Printf("received message: %s", m.Body)
            var code models.Code

            if err := json.Unmarshal(m.Body, &code); err != nil {
                log.Printf("converting to json: %s", err.Error())
                continue
            }

            s.msgHandler.HandleMessage(&code)
        }
    }()

    <-forever
    return nil
}

func (s *RabbitMQReceiver) Close() {
    s.ch.Close()
    s.conn.Close()
}