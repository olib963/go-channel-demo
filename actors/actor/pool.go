package actor

func NewPool[Message any](worker Definition[Message], size int) Actor[Message] {
	workers := make(chan Actor[Message], size)
	for i := 0; i < size; i++ {
		workers <- FromDefinition(worker)
	}
	return FromDefinition(func(ctx Context[Message], message Message) {
		worker := <-workers
		defer func() { workers <- worker }()
		worker.Send(message)
		workers <- worker
	})
}
