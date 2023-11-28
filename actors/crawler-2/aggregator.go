package main

import (
	"github.com/olib963/go-channel-demo/actors/set"
	"github.com/vladopajic/go-actor/actor"
	"log/slog"
	"net/url"
)

type Path = string
type Aggregator struct {
	initialLink   url.URL
	inFlight      set.Set[Path]
	processed     map[Path]set.Set[url.URL]
	MailBox       actor.Mailbox[Parsed]
	WorkerMailbox actor.MailboxSender[Parse]
}

func NewAggregator(WorkerMailbox actor.Mailbox[Parse], initialLink url.URL) *Aggregator {
	return &Aggregator{
		initialLink:   initialLink,
		inFlight:      set.Of(initialLink.Path),
		processed:     make(map[Path]set.Set[url.URL]),
		MailBox:       actor.NewMailbox[Parsed](),
		WorkerMailbox: WorkerMailbox,
	}
}

func (a *Aggregator) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case parsed := <-a.MailBox.ReceiveC():

		a.processed[parsed.Path] = parsed.Urls
		a.inFlight.Remove(parsed.Path)

		for link := range parsed.Urls {
			if a.inFlight.Contains(link.Path) {
				continue
			}

			if _, exists := a.processed[link.Path]; exists {
				continue
			}

			if link.Host == "" {
				link.Host = a.initialLink.Host
			}
			if link.Host != a.initialLink.Host {
				continue
			}

			link.Scheme = a.initialLink.Scheme

			a.inFlight.Add(link.Path)
			slog.Info("New link identified: " + link.String())
			err := a.WorkerMailbox.Send(ctx, Parse{link, a.MailBox})
			if err != nil {
				slog.Error("Failed to send message", "error", err)
				return actor.WorkerEnd
			}
		}
		if len(a.inFlight) > 0 {
			return actor.WorkerContinue
		}
		slog.Info("Finished crawling",
			"processed", len(a.processed),
			"initial_link", a.initialLink.String(),
		)
		return actor.WorkerEnd
	}
}
