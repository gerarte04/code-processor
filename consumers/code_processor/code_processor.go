package code_processor

import (
	rabbMq "http_server/consumers/code_processor/rabbitmq"
	"log"
)

func main() {
    rabbitMqReceiver, err := rabbMq.NewRabbitMQReceiver("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    if err = rabbitMqReceiver.StartReceive(); err != nil {
        rabbitMqReceiver.Close()
        log.Fatalf("%s", err.Error())
    }
}
