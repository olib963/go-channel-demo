package main

import (
	"fmt"
	"github.com/olib963/go-channel-demo/actors/actor"
	"github.com/olib963/go-channel-demo/actors/set"
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
	Reply actor.Actor[Parsed]
}

type Parsed struct {
	Path string
	Urls set.Set[url.URL]
}

func ParseHTML(ctx actor.ActorContext[Parse], toParse Parse) actor.Behaviour {
	client := http.Client{}

	request, err := http.NewRequestWithContext(ctx.Context(), http.MethodGet, toParse.Url.String(), nil)
	if err != nil {
		return actor.Failed(fmt.Errorf("creating request for %s: %w", toParse.Url.String(), err))
	}

	response, err := client.Do(request)
	if err != nil {
		return actor.Failed(fmt.Errorf("fetching %s: %w", toParse.Url.String(), err))
	}

	parsed, err := html.Parse(response.Body)
	if err != nil {
		return actor.Failed(fmt.Errorf("parsing the html of %s: %w", toParse.Url.String(), err))
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
	toParse.Reply.Send(Parsed{toParse.Url.Path, links})
	return actor.Same()
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
