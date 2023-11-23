package actor

import "sync"

func NewPool[Message any](worker Definition[Message], size int) Definition[Message] {
	once := sync.Once{}
	workers := make(chan Actor[Message], size)
	return func(ctx ActorContext[Message], message Message) Behaviour {
		once.Do(func() {
			for i := 0; i < size; i++ {
				workers <- Spawn(ctx, worker)
			}
		})
		worker := <-workers
		defer func() { workers <- worker }()
		worker.Send(message)
		return Same()
	}

}
