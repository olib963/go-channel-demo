package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
	"net/url"
)

func main() {
	initialURL, err := url.Parse("https://monzo.com")
	if err != nil {
		panic(err)
	}

	mainActor := actor.NewDefinition(func(start Start) {
		u := url.URL(start)
		println(u.String())
	})

	mainActor.Send(Start(*initialURL))
}

type Start url.URL
