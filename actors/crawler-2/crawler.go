package main

import (
	"context"
	"github.com/vladopajic/go-actor/actor"
	"net/url"
	"time"
)

func main() {
	initialURL := url.URL{
		Scheme: "https",
		Host:   "monzo.com",
		Path:   "/",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	workerMailbox := actor.NewMailbox[Parse]()

	workers := make([]actor.Actor, 10)
	for i := 0; i < 10; i++ {
		workers[i] = actor.New(&ParserActor{workerMailbox})
	}
	worker := actor.Combine(workers...).Build()

	aggregatorDefinition := NewAggregator(workerMailbox, initialURL)

	toStop := actor.Combine(aggregatorDefinition.MailBox, workerMailbox, worker).Build()

	aggregator := actor.New(aggregatorDefinition, actor.OptOnStop(toStop.Stop))

	a := actor.Combine(aggregator, toStop).
		WithOptions(
			// When every actor stops, stop the system by cancelling the context
			actor.OptOnStopCombined(cancel),
		).Build()

	a.Start()
	// Stop the system if the context is done before the timeout
	defer a.Stop()

	// TODO gotta be a better way.
	err := workerMailbox.Send(ctx, Parse{initialURL, aggregatorDefinition.MailBox})
	if err != nil {
		panic(err)
	}

	<-ctx.Done()

}
