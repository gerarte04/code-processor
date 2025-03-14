package code_processor

type BrokerReceiver interface {
    StartReceive() error
    Close()
}
