package main

import (
	"flag"
	"fmt"
	link "html_link_parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"sitemap_builder/bfs"
	"strings"
)

type flags struct {
	rootUrl  string
	maxDepth uint32
	maxLinks uint32
}

func parseFlags() flags {
	rootUrl := flag.String("root_url", "", "Root Url of the page")
	maxDepth := flag.Uint("max_depth", 3, "Max Depth to search for")
	maxLinks := flag.Uint("max_links", 1000, "Max Links to search for")
	flag.Parse()
	if *rootUrl == "" {
		log.Fatalln("Please provide a root_url, see --help")
	}
	return flags{
		rootUrl:  *rootUrl,
		maxDepth: uint32(*maxDepth),
		maxLinks: uint32(*maxLinks),
	}
}

func hrefs(r io.Reader, base string) (map[string]bool, error) {
	links, err := link.Parse(r)
	if err != nil {
		return nil, err
	}
	ret := map[string]bool{}
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			ret[base+link.Href] = true
		case strings.HasPrefix(link.Href, "/http"):
			ret[link.Href] = true
		}
	}
	return ret, nil
}

func get(urlStr string, keepFn func(string) bool) (map[string]bool, error) {
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	reqUrl := response.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	links, err := hrefs(response.Body, base)
	if err != nil {
		return nil, err
	}
	return filter(links, keepFn), nil
}

func filter(links map[string]bool, keepFn func(string) bool) map[string]bool {
	ret := map[string]bool{}
	for link := range links {
		if keepFn(link) {
			ret[link] = true
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(x string) bool {
		return strings.HasPrefix(x, pfx)
	}
}

func getBaseUrl(urlStr string) (string, error) {
	response, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	reqUrl := response.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	return baseUrl.String(), nil
}

func main() {
	flags := parseFlags()
	baseUrl, err := getBaseUrl(flags.rootUrl)
	if err != nil {
		panic(err)
	}
	keepFn := withPrefix(baseUrl)
	bfsConfig := bfs.NewBFSConfig(
		func(url string) (map[string]bool, error) {
			return get(url, keepFn)
		},
		bfs.WithIgnoreErrors[string](true),
		bfs.WithMaxDepth[string](flags.maxDepth),
		bfs.WithMaxElements[string](uint64(flags.maxLinks)),
	)
	depthNodes, err := bfsConfig.BFS(flags.rootUrl)
	if err != nil {
		panic(err)
	}
	for _, depthNode := range depthNodes {
		fmt.Printf("url: %s\tdepth: %d\n", depthNode.Node, depthNode.Depth)
	}
}
