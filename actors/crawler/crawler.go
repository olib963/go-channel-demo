package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
	"net/url"
	"time"
)

func main() {
	actor.NewSystem(Crawler).
		WithTimeout(5 * time.Minute).
		Start().
		AwaitTermination()
}

func Crawler(initialContext actor.Context) {
	initialURL := url.URL{
		Scheme: "https",
		Host:   "monzo.com",
		Path:   "/",
	}

	workerDefinition := ParseHTML
	pool := actor.NewPool(workerDefinition, 10)
	workers := actor.Spawn(initialContext, pool)
	// TODO better state management, definition should encapsulate state
	agg := actor.Spawn(initialContext, aggregator(initialURL, workers))
	workers.Send(Parse{initialURL, agg})
}
