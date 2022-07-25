package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func init() {
	defaultChapterHandlerTmpl = template.Must(template.New("").Parse(defaultChapterHandlerTmplString))
}

var defaultChapterHandlerTmpl *template.Template

const defaultChapterHandlerTmplString string = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`

func defaultChapterParser(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro" // default to /intro
	}
	return path[1:] // slicing off the leading slash "/"
}

type storyHandlerOpts func(*storyHandler)

func withTemplate(tmpl *template.Template) storyHandlerOpts {
	return func(sh *storyHandler) {
		sh.tmpl = tmpl
	}
}

func withChapterParser(chapterParser func(r *http.Request) string) storyHandlerOpts {
	return func(sh *storyHandler) {
		sh.chapterParser = chapterParser
	}
}

type storyHandler struct {
	s             story
	tmpl          *template.Template
	chapterParser func(r *http.Request) string
}

func mkStoryHandler(s story, opts ...storyHandlerOpts) http.Handler {
	if s == nil {
		s = make(map[string]chapter)
	}
	sh := storyHandler{s, defaultChapterHandlerTmpl, defaultChapterParser}
	for _, opt := range opts {
		opt(&sh)
	}
	return sh
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chapterKey := h.chapterParser(r)

	if chapter, ok := h.s[chapterKey]; ok {
		err := h.tmpl.Execute(w, chapter)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Chapter not found.", http.StatusNotFound)
	}
}
