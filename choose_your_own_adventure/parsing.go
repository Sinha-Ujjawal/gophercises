package main

import (
	"encoding/json"
	"io/ioutil"
)

type option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type story = map[string]chapter

func readStoryFromJSON(jsonPath string) (story, error) {
	bytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	var story map[string]chapter
	json.Unmarshal(bytes, &story)
	return story, nil
}
