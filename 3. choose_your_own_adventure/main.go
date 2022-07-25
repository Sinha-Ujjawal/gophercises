package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type flags struct {
	storyJSONPath string
	port          uint16
}

func parseFlags() flags {
	storyJSONPath := flag.String("storyJson", "./gopher.json", "JSON Path of the story data")
	port := flag.Int("port", 3030, "Port of the web server")
	flag.Parse()
	return flags{storyJSONPath: *storyJSONPath, port: uint16(*port)}
}

func main() {
	flags := parseFlags()
	story, err := readStoryFromJSON(flags.storyJSONPath)
	if err != nil {
		log.Fatal("Could not parse story from the json path!\n", err)
	}
	storyHandler := mkStoryHandler(story)
	port := fmt.Sprintf(":%d", flags.port)
	log.Printf("Starting the server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, storyHandler))
}
