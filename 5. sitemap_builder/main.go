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

const MAX_DEPTH = 3

type flags struct {
	rootUrl  string
	maxDepth uint32
}

func parseFlags() flags {
	rootUrl := flag.String("root_url", "", "Root Url of the page")
	maxDepth := flag.Uint("max_depth", MAX_DEPTH, "Max Depth to search for")
	flag.Parse()
	if *rootUrl == "" {
		log.Fatalln("Please provide a root_url, see --help")
	}
	if *maxDepth > MAX_DEPTH {
		fmt.Printf("--max_depth should be between 0-%d. The script would run with max_depth: %d\n", MAX_DEPTH, MAX_DEPTH)
		*maxDepth = MAX_DEPTH
	}
	return flags{
		rootUrl:  *rootUrl,
		maxDepth: uint32(*maxDepth),
	}
}

func hrefs(r io.Reader, base string) ([]string, error) {
	links, err := link.Parse(r)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			ret = append(ret, base+link.Href)
		case strings.HasPrefix(link.Href, "/http"):
			ret = append(ret, link.Href)
		}
	}
	return ret, nil
}

func get(urlStr string) ([]string, error) {
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
	return links, nil
}

func main() {
	flags := parseFlags()
	fmt.Println("Running the script with:")
	fmt.Printf("  --root_url: %s\n", flags.rootUrl)
	fmt.Printf("  --max_depth: %d\n", flags.maxDepth)
	fmt.Println()
	for depthNode := range bfs.BFS(flags.rootUrl, get, flags.maxDepth) {
		if depthNode.Err != nil {
			errorMessage := fmt.Sprintf("Error(%s) caught when processing: %s", depthNode.Err, depthNode.Node)
			fmt.Println(errorMessage)
			panic(errorMessage)
		}
		pad := strings.Repeat(" ", int(depthNode.Depth)*4)
		fmt.Printf("%s%s\n", pad, depthNode.Node)
	}
}
