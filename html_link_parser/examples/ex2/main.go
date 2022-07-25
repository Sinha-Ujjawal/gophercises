package main

import (
	"fmt"
	link "html_link_parser"
	"strings"
)

const exampleHTML string = `
<a href="#">
  Something here <a href="/dog">nested dog link</a>
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
