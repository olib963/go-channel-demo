package main

import (
	"fmt"
	"sync"
	"time"
)

func LimiterExample() {
	pushBased := NewPushLimiter(3)
	for i := 0; i < 5; i++ {
		GoPush(pushBased, func() {
			fmt.Println("Hello World")
			time.Sleep(1 * time.Second)
		})
	}
	pushBased.Wait()

	time.Sleep(1 * time.Second)

	pullBased := NewPullLimiter(3, 1, func(i int) int { return i + 1 })
	for i := 0; i < 10; i++ {
		GoPull(pullBased, func(threadID int) {
			fmt.Printf("Running from thread %d\n", threadID)
			time.Sleep(time.Duration(threadID) * time.Second)
		})
	}
	pullBased.Wait()

}

type pushLimiter struct {
	c chan struct{}
	*sync.WaitGroup
}

// Either create a channel of empty structs and write to it as the limiter

func NewPushLimiter(size uint) pushLimiter {
	return pushLimiter{make(chan struct{}, size), &sync.WaitGroup{}}
}

func GoPush(p pushLimiter, f func()) {
	go func() {
		p.Add(1)
		p.c <- struct{}{}
		defer func() {
			<-p.c
			p.Done()
		}()
		f()
	}()
}

type pullLimiter[A any] struct {
	c chan A
	*sync.WaitGroup
}

// Or create a channel full of "tokens" to pull from as the limiter

func NewPullLimiter[A any](size uint, initial A, next func(A) A) pullLimiter[A] {
	as := make([]A, size)
	current := initial
	as[0] = current
	for i := uint(1); i < size; i++ {
		as[i] = next(current)
	}
	return NewPullLimiterFrom(as)
}

func NewPullLimiterFrom[A any](list []A) pullLimiter[A] {
	as := make(chan A, len(list))
	for _, a := range list {
		as <- a
	}
	return pullLimiter[A]{as, &sync.WaitGroup{}}
}

func GoPull[A any](p pullLimiter[A], f func(A)) {
	go func() {
		p.Add(1)
		a := <-p.c
		defer func() {
			p.c <- a // Put the token back when done
			p.Done()
		}()
		f(a)
	}()
}
