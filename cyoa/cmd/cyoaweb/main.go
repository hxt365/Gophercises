package main

import (
	"Gophercises/cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filename := flag.String("filename", "gopher.json", "The json story file")
	port := flag.Int("port", 3000, "The port number for CYOA web")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot open the json file %s: %s", *filename, err))
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse the json file %s: %s", *filename, err))
	}

	//tmpl := template.Must(template.New("").Parse(storyTmpl))
	//h := cyoa.NewHandler(story,
	//	cyoa.WithTemplate(tmpl),
	//	cyoa.WithPathFn(pathFn),
	//)
	h := cyoa.NewHandler(story)
	mux := http.NewServeMux()
	mux.Handle("/", h)

	fmt.Printf("Starting server on port %v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), mux))
}

var storyTmpl = `
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

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
