package main

import (
	"fmt"
	link "html_link_parser"
	"strings"
)

const exampleHTML string = `
<a href="/dog">
  Text not in a span
  <span>Something in a span</span>
  <b>Bold text!</b>
</a>
`

func main() {
	r := strings.NewReader(exampleHTML)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
