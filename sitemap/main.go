package main

import (
	"Gophercises/link"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	url := flag.String("url", "https://gophercises.com", "A website url that you want to build a sitemap for")
	maxDepth := flag.Int("maxDepth", 5, "Link depth")
	flag.Parse()

	links, err := bfs(*url, *maxDepth)
	if err != nil {
		panic(err)
	}

	if err = toXML(links); err != nil {
		panic(err)
	}
}

func bfs(url string, maxDepth int) ([]string, error) {
	var q, res []string
	nq := []string{url}
	visited := make(map[string]struct{})

	for i := 0; i < maxDepth; i++ {
		q, nq = nq, []string{}
		if len(q) == 0 {
			break
		}
		for _, l := range q {
			if _, ok := visited[l]; ok {
				continue
			}
			res = append(res, l)
			subLinks, err := get(l)
			if err != nil {
				return nil, err
			}
			visited[l] = struct{}{}
			nq = append(nq, subLinks...)
		}
	}
	return res, nil
}

func get(urlStr string) ([]string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Host:   reqUrl.Host,
		Scheme: reqUrl.Scheme,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base)), nil
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var res []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "http"):
			res = append(res, strings.TrimRight(l.Href, "/"))
		case strings.HasPrefix(l.Href, "/"):
			res = append(res, strings.TrimRight(base+l.Href, "/"))
		}
	}
	return res
}

func filter(links []string, keepFn func(string) bool) []string {
	var res []string
	for _, h := range links {
		if keepFn(h) {
			res = append(res, h)
		}
	}
	return res
}

func withPrefix(pfx string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, pfx)
	}
}

func toXML(links []string) error {
	fmt.Print(xml.Header)
	XML := urlset{
		Xmlns: xmlns,
	}
	for _, l := range links {
		XML.Urls = append(XML.Urls, loc{l})
	}
	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")
	err := encoder.Encode(XML)
	fmt.Println()
	return err
}
