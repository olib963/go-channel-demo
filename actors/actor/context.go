package actor

type Context[Message any] interface {
	Self() Actor[Message]
}

type actorContext[Message any] struct {
	self Actor[Message]
}

func (c actorContext[Message]) Self() Actor[Message] { return c.self }
