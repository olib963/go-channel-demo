package actor

import (
	"context"
	"log/slog"
)

type Actor[Message any] interface {
	Send(Message)
	Stop()
}

type Definition[Message any] func(ActorContext[Message], Message) Behaviour

func Spawn[Message any](ctx Context, name string, definition Definition[Message]) Actor[Message] {
	messages := make(chan Message)

	slog.Info("Spawned actor", "name", name)

	ctx.actorGroup().Add(1)

	internalCtx, cancel := context.WithCancel(ctx.Context())

	a := actor[Message]{messages, cancel}
	newContext := actorContext[Message]{
		self:        a,
		internalCtx: internalCtx,
		actors:      ctx.actorGroup(),
	}

	go func() {
		defer ctx.actorGroup().Done()
		for {
			select {
			// Await termination.
			case <-newContext.Context().Done():
				slog.Info("Terminated actor", "name", name)
				return
			case message := <-messages:
				b := definition(newContext, message)
				switch b := b.(type) {
				case stop:
					slog.Info("Stopping actor", "name", name)
					cancel()
				case failed:
					slog.Error("Actor failed", "name", name, "error", b.err)
					panic(b.err)
				case same:
					// Keep processing.
				default:
					panic("unknown behaviour. This is a bug in the actor library")
				}
			}
		}
	}()
	return a
}

type actor[Message any] struct {
	messages chan Message
	stop     context.CancelFunc
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

func (a actor[Message]) Stop() { a.stop() }
