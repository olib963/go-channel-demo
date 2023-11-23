package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
	"net/url"
	"time"
)

func main() {
	actor.NewSystem(Crawler).
		WithTimeout(5 * time.Minute).
		Start()
}

func Crawler(initialContext actor.Context[struct{}]) {
	initialURL := url.URL{
		Scheme: "https",
		Host:   "monzo.com",
		Path:   "/",
	}

	workerDefinition := ParseHTML
	pool := actor.NewPool(workerDefinition, 10)
	agg := actor.FromDefinition(aggregator(initialURL, pool))
	pool.Send(Parse{initialURL, agg})
}
