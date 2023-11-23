package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
	"github.com/olib963/go-channel-demo/actors/set"
	"log/slog"
	"net/url"
)

type Path = string

func aggregator(initialLink url.URL, worker actor.Actor[Parse]) actor.Definition[Parsed] {
	inFlight := set.Of(initialLink.Path)
	processed := make(map[Path]set.Set[url.URL])

	return func(ctx actor.ActorContext[Parsed], parsed Parsed) actor.Behaviour {
		processed[parsed.Path] = parsed.Urls
		inFlight.Remove(parsed.Path)

		for link := range parsed.Urls {
			if inFlight.Contains(link.Path) {
				continue
			}

			if _, exists := processed[link.Path]; exists {
				continue
			}

			if link.Host == "" {
				link.Host = initialLink.Host
			}
			if link.Host != initialLink.Host {
				continue
			}

			link.Scheme = initialLink.Scheme

			inFlight.Add(link.Path)
			slog.Info("New link identified: " + link.String())
			worker.Send(Parse{link, ctx.Self()})
		}
		if len(inFlight) > 0 {
			return actor.Same()
		}
		slog.Info("Finished crawling",
			"processed", len(processed),
			"initial_link", initialLink.String(),
		)
		worker.Stop()
		return actor.Stop()

	}
}
