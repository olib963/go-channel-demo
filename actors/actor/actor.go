package actor

import "context"

type Actor[Message any] interface {
	Send(Message)
}

type Definition[Message any] func(Context[Message], Message) Behaviour

func Spawn[Any, Message any](ctx Context[Any], fn Definition[Message]) Actor[Message] {
	messages := make(chan Message, 1)

	a := actor[Message]{messages}

	internalCtx, cancel := context.WithCancel(ctx.ctx())
	newContext := actorContext[Message]{
		self:        a,
		internalCtx: internalCtx,
	}

	go func() {
		for {
			select {
			// Await termination.
			case <-newContext.ctx().Done():
				return
			case message := <-messages:
				b := fn(newContext, message)
				switch b := b.(type) {
				case stop:
					cancel()
				case failed:
					// TODO propagate error up ctx chain.
					panic(b.err)
				case same:
					// Keep processing.
				default:
					panic("unknown behaviour")
				}
			}
		}
	}()
	return a
}

type actor[Message any] struct {
	messages chan Message
}

func (a actor[Message]) Send(message Message) {
	go func() {
		a.messages <- message
	}()
}
