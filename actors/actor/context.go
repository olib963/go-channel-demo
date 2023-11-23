package actor

import "context"

type Context[Message any] interface {
	Self() Actor[Message]
	ctx() context.Context
}

type actorContext[Message any] struct {
	self        Actor[Message]
	internalCtx context.Context
}

func (c actorContext[Message]) Self() Actor[Message] { return c.self }

func (c actorContext[Message]) ctx() context.Context {
	return c.internalCtx
}
