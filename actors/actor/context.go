package actor

import (
	"context"
	"sync"
)

type Context interface {
	Context() context.Context
	actorGroup() *sync.WaitGroup
}

type ActorContext[Message any] interface {
	Context
	Self() Actor[Message]
}

type actorContext[Message any] struct {
	self        Actor[Message]
	internalCtx context.Context
	actors      *sync.WaitGroup
}

func (c actorContext[Message]) Self() Actor[Message] { return c.self }

func (c actorContext[Message]) Context() context.Context { return c.internalCtx }

func (c actorContext[Message]) actorGroup() *sync.WaitGroup { return c.actors }
