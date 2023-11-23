package main

import (
	"github.com/olib963/go-channel-demo/actors/actor"
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
	Urls set[url.URL]
}

func ParseHTML(toParse Parse) {
	response, err := http.Get(toParse.Url.String())
	if err != nil {
		slog.Info("Error fetching %s: %s", toParse.Url.String(), err.Error())
		return
	}

	parsed, err := html.Parse(response.Body)
	if err != nil {
		slog.Info("Error parsing %s: %s", toParse.Url.String(), err.Error())
		return
	}

	links := make([]url.URL, 0)
	for _, anchor := range allAnchors(parsed) {
		href, exists := findHref(anchor)
		if !exists {
			slog.Warn("Anchor has no href", anchor)
			continue
		}
		url, err := url.Parse(href)
		if err != nil {
			slog.Warn("Error parsing URL", href, err)
			continue
		}
		links = append(links, *url)
	}
	slog.Info("Found links", links)
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
