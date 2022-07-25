package main

import (
	"database/sql"
	"fmt"
	"net/http"
	url_shortner "url_shortner_app/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	mux := defaultMux()

	pathsToUrl := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := url_shortner.MapHandler(pathsToUrl, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := url_shortner.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := `
[
	{"path": "/google", "url": "https://www.google.com"}
]
	`

	jsonHandler, err := url_shortner.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		panic(err)
	}

	dbHandler := url_shortner.DBHandler(db, "url_shortner", "path", "dest", jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}
