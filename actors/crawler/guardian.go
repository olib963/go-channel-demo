package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
	"log/slog"
	"net/url"
)

func aggregator(initialLink url.URL, worker actor.Actor[Parse]) actor.Definition[Parsed] {
	inFlight := set[string]{initialLink.Path: struct{}{}}
	processed := make(map[string]set[url.URL])

	return func(ctx actor.Context[Parsed], parsed Parsed) {
		processed[parsed.Path] = parsed.Urls
		delete(inFlight, parsed.Path)

		for link := range parsed.Urls {
			if _, exists := processed[link.Path]; exists {
				continue
			}
			if _, exists := inFlight[link.Path]; exists {
				continue
			}
			if link.Host == "" {
				link.Host = initialLink.Host
			}
			if link.Host != initialLink.Host {
				continue
			}

			link.Scheme = initialLink.Scheme

			inFlight[link.Path] = struct{}{}
			slog.Info("Crawling %s", link.String())
			worker.Send(Parse{link, ctx.Self()})
		}
		if len(inFlight) == 0 {
			// TODO respond
			slog.Info("Finished crawling")
		}
	}
}

type set[T comparable] map[T]struct{}
