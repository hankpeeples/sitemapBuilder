package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

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

	links, _ := link.Parse(resp.Body)

	for _, l := range links {
		fmt.Println(l)
	}
}
