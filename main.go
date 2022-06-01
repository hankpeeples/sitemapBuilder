package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/hankpeeples/linkParser"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the max number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	// pages := get(*urlFlag)

	fmt.Printf("Found pages: ⬇️\n\n")
	for _, page := range pages {
		fmt.Println(page)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	// keep track of all urls that have been visited
	// https://dave.cheney.net/2014/03/25/the-empty-struct
	seen := make(map[string]struct{})
	// string is the key, struct is the type
	var q map[string]struct{}
	// all unseen links, these are child links from links that have been visited
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		// `q` will become whatever is in `nq`, then make a new `nq` for use
		q, nq = nq, make(map[string]struct{})
		for page, _ := range q {
			// if `page` was `seen` in the map, `ok` will be true
			if _, ok := seen[page]; ok {
				// skip it
				continue
			}
			// mark page as seen
			seen[page] = struct{}{}
			for _, l := range get(page) {
				// put each link in the next queue
				nq[l] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for page, _ := range seen {
		ret = append(ret, page)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
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

	// fmt.Println("✅ Request made to:", urlStr)

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

	// fmt.Printf("⎵ Base URL: %s\n\n", base)

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)

	var ret []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, l := range links {
		// only keep links that have the same domain as the original (base)
		if keepFn(l) {
			ret = append(ret, l)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(l string) bool {
		return strings.HasPrefix(l, pfx)
	}
}
