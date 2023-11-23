package actor

type Actor[Message any] interface {
	Send(Message)
}

type Definition[Message any] func(Context[Message], Message) // TODO return behaviour.

// TODO _spawn_ on Context
func FromDefinition[Message any](fn Definition[Message]) Actor[Message] {
	return actor[Message](fn)
}

type actor[Message any] func(Context[Message], Message)

func (a actor[Message]) Send(message Message) {
	context := actorContext[Message]{self: a}
	a(context, message)
}
