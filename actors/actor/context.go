package actor

import (
	"context"
)

type Context interface {
	Context() context.Context
}

type ActorContext[Message any] interface {
	Context
	Self() Actor[Message]
}

type actorContext[Message any] struct {
	self        Actor[Message]
	internalCtx context.Context
}

func (c actorContext[Message]) Self() Actor[Message] { return c.self }

func (c actorContext[Message]) Context() context.Context { return c.internalCtx }
