package main

import (
	"Gophercises/link"
	"fmt"
	"log"
	"strings"
)

var htmlContent = `
<html>
<head>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
<h1>Social stuffs</h1>
<div>
    <a href="https://www.twitter.com/joncalhoun">
        Check me out on twitter
        <i class="fa fa-twitter" aria-hidden="true"></i>
    </a>
    <a href="https://github.com/gophercises">
        Gophercises is on <strong>Github</strong>!
    </a>
</div>
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
