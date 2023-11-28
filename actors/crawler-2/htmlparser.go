package main

import (
	"fmt"
	"github.com/olib963/go-channel-demo/actors/set"
	"github.com/vladopajic/go-actor/actor"
	"golang.org/x/net/html"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	anchor = "a"
	href   = "href"
)

type Parse struct {
	Url   url.URL
	Reply actor.MailboxSender[Parsed]
}

type Parsed struct {
	Path string
	Urls set.Set[url.URL]
}

type ParserActor struct {
	Mailbox actor.Mailbox[Parse]
}

func (a *ParserActor) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case toParse := <-a.Mailbox.ReceiveC():
		client := http.Client{}

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, toParse.Url.String(), nil)
		if err != nil {
			slog.Error(fmt.Sprintf("creating request for %s", toParse.Url.String()), "error", err)
			return actor.WorkerEnd
		}

		response, err := client.Do(request)
		if err != nil {
			slog.Error(fmt.Sprintf("fetching %s", toParse.Url.String()), "error", err)
			return actor.WorkerEnd
		}

		parsed, err := html.Parse(response.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("parsing the html of %s", toParse.Url.String()), "error", err)
			return actor.WorkerEnd
		}

		links := set.Set[url.URL]{}
		for _, anchor := range allAnchors(parsed) {
			href, exists := findHref(anchor)
			if !exists {
				continue
			}
			url, err := url.Parse(href)
			if err != nil {
				slog.Warn("Error parsing URL", "href", href, "err", err)
				continue
			}
			links[*url] = struct{}{}
		}

		if err := toParse.Reply.Send(ctx, Parsed{toParse.Url.Path, links}); err != nil {
			slog.Error(fmt.Sprintf("sending parsed links for %s", toParse.Url.String()), "error", err)
			return actor.WorkerEnd
		}

		return actor.WorkerContinue
	}
}

func allAnchors(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == anchor {
		return []*html.Node{node}
	}

	nodes := make([]*html.Node, 0)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, allAnchors(c)...)
	}
	return nodes
}

func findHref(node *html.Node) (string, bool) {
	for _, attribute := range node.Attr {
		if attribute.Key == href {
			return attribute.Val, true
		}
	}
	return "", false
}
