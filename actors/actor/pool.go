package actor

import (
	"strconv"
	"sync"
)

func NewPool[Message any](worker Definition[Message], poolName string, size int) Definition[Message] {
	once := sync.Once{}
	workers := make(chan Actor[Message], size)
	return func(ctx ActorContext[Message], message Message) Behaviour {
		once.Do(func() {
			for i := 0; i < size; i++ {
				name := poolName + "-" + strconv.Itoa(i)
				workers <- Spawn(ctx, name, worker)
			}
		})
		worker := <-workers
		defer func() { workers <- worker }()
		worker.Send(message)
		return Same()
	}

}
