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
	AwaitTermination()
}

type system struct {
	setup  func(ctx Context)
	ctx    context.Context
	cancel context.CancelFunc
}

func (s system) WithTimeout(timeout time.Duration) System {
	newContext, cancel := context.WithTimeout(s.ctx, timeout)
	return system{s.setup, newContext, cancel}
}

func (s system) Start() RunningSystem {
	s.setup(actorContext[struct{}]{
		internalCtx: s.ctx,
	})
	return s
}

func (s system) StopNow() {
	s.cancel()
}

func (s system) AwaitTermination() {
	<-s.ctx.Done()
}

func NewSystem(setup func(ctx Context)) System {
	ctx, cancel := context.WithCancel(context.Background())
	return system{setup, ctx, cancel}
}
