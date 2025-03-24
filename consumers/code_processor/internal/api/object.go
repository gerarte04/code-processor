package api

type MessageHandler interface {
    HandleMessage(message []byte)
}
