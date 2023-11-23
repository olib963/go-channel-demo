package actor

import (
	"context"
	"time"
)

type System interface {
	WithTimeout(timeout time.Duration) System
	Start() RunningSystem
}

type RunningSystem interface {
	StopNow()
}

type system struct {
	setup  func(ctx Context[struct{}])
	ctx    context.Context
	cancel context.CancelFunc
}

func (s system) WithTimeout(timeout time.Duration) System {
	newContext, cancel := context.WithTimeout(s.ctx, timeout)
	return system{s.setup, newContext, cancel}
}

func (s system) Start() RunningSystem {
	initialActorContext := actorContext[struct{}]{
		self: FromDefinition(func(ctx Context[struct{}], message struct{}) {}),
	}
	s.setup(initialActorContext)
	return s
}

func (s system) StopNow() {
	s.cancel()
}

func NewSystem(setup func(ctx Context[struct{}])) System {
	ctx, cancel := context.WithCancel(context.Background())
	return system{setup, ctx, cancel}
}
