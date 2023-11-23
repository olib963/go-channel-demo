package actor

import "context"

type Actor[Message any] interface {
	Send(Message)
}

type Definition[Message any] func(ActorContext[Message], Message) Behaviour

func Spawn[Message any](ctx Context, fn Definition[Message]) Actor[Message] {
	messages := make(chan Message)

	a := actor[Message]{messages}

	internalCtx, cancel := context.WithCancel(ctx.Context())
	newContext := actorContext[Message]{
		self:        a,
		internalCtx: internalCtx,
	}

	go func() {
		for {
			select {
			// Await termination.
			case <-newContext.Context().Done():
				return
			case message := <-messages:
				b := fn(newContext, message)
				switch b := b.(type) {
				case stop:
					cancel()
				case failed:
					// TODO propagate error up Context chain.
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
	// Spawning a goroutine here avoids deadlocks at the cost of infinite
	// goroutine creation (and thus memory usage).
	// We could more strictly control the actor concurrency, but this is more
	// difficult _and_ a tradeoff we're willing to make. Deadlocks are more
	// likely to be a problem than memory usage.
	go func() {
		a.messages <- message
	}()
}
