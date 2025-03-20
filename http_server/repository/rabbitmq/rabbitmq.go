package rabbitmq

import (
	"encoding/json"
	"fmt"
	"http_server/config"
	"http_server/repository/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQSender struct {
    conn *amqp.Connection
    ch *amqp.Channel
    queue *amqp.Queue

    cfg config.RabbitMQConfig
}

func NewRabbitMQSender(cfg config.RabbitMQConfig) (*RabbitMQSender, error) {
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

    return &RabbitMQSender{
        conn: conn,
        ch: ch,
        queue: &queue,
        cfg: cfg,
    }, nil
}

func (s *RabbitMQSender) Send(message *models.Code) error {
    data, err := json.Marshal(CreateRabbitMQMessage(message))

    if err != nil {
        return fmt.Errorf("formatting rabbitmq message: %s", err.Error())
    }

    err = s.ch.Publish(
        "",
        s.queue.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body: data,
        },
    )

    if err != nil {
        return fmt.Errorf("sending rabbitmq message: %s", err.Error())
    }

    return nil
}

func (s *RabbitMQSender) Close() {
    s.ch.Close()
    s.conn.Close()
}
