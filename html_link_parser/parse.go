package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represent a link (<a href"...">) in an HTML
// document.
type Link struct {
	Href string
	Text string
}

// Parse will take an HTML document and will return a slice
// of links parsed from it
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	links := make([]Link, 0)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(node *html.Node) Link {
	var ret Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	var sb strings.Builder
	text(node, &sb)
	ret.Text = sb.String()
	return ret
}

func text(node *html.Node, sb *strings.Builder) {
	if node.Type == html.TextNode {
		sb.WriteString(" ")
		sb.WriteString(strings.Join(strings.Fields(node.Data), " "))
		return
	}
	if node.Type != html.ElementNode {
		sb.WriteString("")
		return
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text(c, sb)
	}
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
