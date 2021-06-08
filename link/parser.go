package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	linkNodes := getLinkNodes(doc)
	links := buildLinks(linkNodes)
	return links, nil
}

func getLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var linkNodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		linkNodes = append(linkNodes, getLinkNodes(c)...)
	}
	return linkNodes
}

func buildLinks(linkNodes []*html.Node) []Link {
	var links []Link
	for _, n := range linkNodes {
		links = append(links, Link{
			Href: getHref(n),
			Text: getText(n),
		})
	}
	return links
}

func getHref(linkNode *html.Node) string {
	for _, a := range linkNode.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func getText(linkNode *html.Node) string {
	if linkNode.Type == html.TextNode {
		return strings.TrimSpace(linkNode.Data)
	}
	text := ""
	for c := linkNode.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return text
}
