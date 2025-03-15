package internal

type BrokerReceiver interface {
    StartReceive() error
    Close()
}
