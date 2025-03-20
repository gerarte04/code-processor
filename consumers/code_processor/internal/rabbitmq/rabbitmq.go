package rabbitmq

import (
	"code_processor/config"
	"code_processor/internal/api"
	"code_processor/internal/models"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQReceiver struct {
    conn *amqp.Connection
    ch *amqp.Channel
    queue *amqp.Queue

    cfg config.RabbitMQConfig
    msgHandler api.MessageHandler
}

func NewRabbitMQReceiver(cfg config.RabbitMQConfig, msgHandler api.MessageHandler) (*RabbitMQReceiver, error) {
    url := fmt.Sprintf("amqp://%s@%s:%s", cfg.Authority, cfg.Host, cfg.Port)
    conn, err := amqp.Dial(url)

    if err != nil {
        return nil, fmt.Errorf("connecting rabbitmq: %s", err.Error())
    }

    ch, err := conn.Channel()

    if err != nil {
        return nil, fmt.Errorf("opening rabbitmq channel: %s", err.Error())
    }

    queue, err := ch.QueueDeclare(
        cfg.QueueName,
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

        cfg: cfg,
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