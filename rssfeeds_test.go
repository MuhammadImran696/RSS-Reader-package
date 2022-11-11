package rssfeeds

import (
	"fmt"
	"testing"
)

func TestParseMethod(t *testing.T) {
	urls := []string{"http://feeds.nytimes.com/nyt/rss/HomePage", "https://www.latimes.com/world/rss2.0.xml"}
	// var data []RssItem
	data := Parse(urls)
	if len(data) == 0 {
		t.Fatalf(`No data Found`)
	}
}
func TestParseWithWrongURL(t *testing.T) {
	urls := []string{"http://feeds.nytimes.com/"}

	data := Parse(urls)
	fmt.Print(data)
	if len(data) > 0 {
		t.Errorf("got %q, wanted %v", data, nil)
	}
}

func TestRequestMethod(t *testing.T) {
	urls := "http://feeds.nytimes.com/nyt/rss/HomePage"
	// var data []RssItem
	data := doRequest(urls)
	if len(data) == 0 {
		t.Fatalf(`No data Found`)
	}
}
