package actor

type Actor[Message any] interface {
	Send(Message)
}

func NewDefinition[Message any](fn func(Message)) Actor[Message] {
	return actor[Message](fn)
}

type actor[Message any] func(Message)

func (a actor[Message]) Send(message Message) {
	a(message)
}
