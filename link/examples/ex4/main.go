package main

import (
	"Gophercises/link"
	"fmt"
	"log"
	"strings"
)

var htmlContent = `
<html>
<body>
<a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>`

func main() {
	r := strings.NewReader(htmlContent)
	links, err := link.Parse(r)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse the html: %s", err))
	}
	for i, l := range links {
		fmt.Printf("Link %d:\n   Href: %s\n   Text: %s\n", i, l.Href, l.Text)
	}
}
