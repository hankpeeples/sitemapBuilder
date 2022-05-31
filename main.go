package main

import (
	"flag"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/hankpeeples/linkParser"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	// close resp.Body after function has finished (in this case, closes after main finishes)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// final URL after any redirects if there happen to be any
	reqUrl := resp.Request.URL
	// base url is the base entry point of the site
	// Ex. If `reqUrl = https://gophercises.com/demos/cyoa`
	// 		The base url is now `https://gophercises.com`
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	links, _ := link.Parse(resp.Body)

	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
}
