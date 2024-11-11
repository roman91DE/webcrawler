package main

import (
	"fmt"
	"github.com/roman91DE/webcrawler/pkg/crawler"
	"os"
)

func main() {
	urls := os.Args[1:]
	crawl(urls)
}

func crawl(urls []string) {
	fetched := crawler.FetchAll(urls)
	links := make([]string, 0)

	for _, hit := range fetched {
		if hit.Error != nil {
			continue
		}
		linksRaw, err := crawler.ExtractLinks(&hit)
		if err != nil {
			continue
		}
		links = append(links, crawler.FilterLinks(linksRaw)...)
	}
	for _, link := range links {
		fmt.Println(link)
	}
	crawl(links)
}
