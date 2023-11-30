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

	pool := actor.NewPool(ParseHTML, "html-parser", 10)
	workers := actor.Spawn(initialContext, "worker-pool", pool)
	agg := actor.Spawn(initialContext, "aggregator", aggregator(initialURL, workers))
	workers.Send(Parse{initialURL, agg})
}
