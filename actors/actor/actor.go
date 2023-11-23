package actor

type Actor[Message any] interface {
	Send(Message)
}

type Definition[Message any] func(Context[Message], Message) // TODO return behaviour.

func Spawn[Any, Message any](ctx Context[Any], fn Definition[Message]) Actor[Message] {
	messages := make(chan Message, 1)

	a := actor[Message]{messages}
	newContext := actorContext[Message]{
		self:        a,
		internalCtx: ctx.ctx(),
	}

	// TODO handle behaviours and allow immediate returns rather than blocking.
	go func() {
		for {
			select {
			case <-newContext.ctx().Done():
				return
			case message := <-messages:
				fn(newContext, message)
			}
		}
	}()
	return a
}

type actor[Message any] struct {
	messages chan Message
}

func (a actor[Message]) Send(message Message) {
	a.messages <- message
}
