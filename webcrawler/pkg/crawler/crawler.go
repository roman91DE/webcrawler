package crawler

import (
	"golang.org/x/net/html"
	"net/url"
	"regexp"
	"strings"
)

func ExtractLinks(f *FetchResult) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(f.Content))
	if err != nil {
		return nil, err
	}

	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if _, err := url.Parse(attr.Val); err == nil {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return links, nil
}

func FilterLinks(links []string) []string {

	buf := make(map[string]bool)
	validURL := regexp.MustCompile(`^https?://`)

	for _, url := range links {
		if validURL.MatchString(url) {
			buf[url] = true
		}
	}

	results := make([]string, 0, len(buf))
	for k := range buf {
		results = append(results, k)
	}

	return results
}
