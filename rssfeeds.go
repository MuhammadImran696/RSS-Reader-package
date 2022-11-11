package rssfeeds

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Dc      string   `xml:"dc,attr"`
	Content string   `xml:"content,attr"`
	Atom    string   `xml:"atom,attr"`
	Media   string   `xml:"media,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Image struct {
			Text  string `xml:",chardata"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Generator     string `xml:"generator"`
		LastBuildDate string `xml:"lastBuildDate"`
		Copyright     string `xml:"copyright"`
		Language      string `xml:"language"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
			Guid        struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Creator string `xml:"creator"`
			PubDate string `xml:"pubDate"`
			Content struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Medium string `xml:"medium,attr"`
				Height string `xml:"height,attr"`
				Width  string `xml:"width,attr"`
			} `xml:"content"`
		} `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string    `xml:"title"`
	Source      string    `xml:"source"`
	SourceURL   string    `xml:"source_url"`
	Link        string    `xml:"link"`
	PublishDate time.Time `xml:"publish_date"`
	Description string    `xml:"description"`
}

func Parse(urls []string) []RssItem {

	var wg sync.WaitGroup
	var content []RssItem
	var collection []RssItem
	var mut sync.Mutex
	for _, url := range urls {

		wg.Add(1)
		go func(url string) {
			mut.Lock()
			content = doRequest(url)
			collection = append(collection, content...)
			mut.Unlock()
			wg.Done()
		}(url)
	}
	wg.Wait()
	return collection
}

func doRequest(url string) (content []RssItem) {

	data, err := http.Get(url)

	if err != nil {

		log.Println(err)
		return
	}

	defer data.Body.Close()
	var RssFeeds rss
	body, err := ioutil.ReadAll(data.Body)

	if err != nil {

		log.Println(err)
		return
	}
	stringbody := string(body)
	// fmt.Print(stringbody)
	data.Body.Close()

	xml.Unmarshal([]byte(stringbody), &RssFeeds)
	items := []RssItem{}
	item := RssItem{}
	for i := range RssFeeds.Channel.Item {
		item.Title = RssFeeds.Channel.Item[i].Title
		item.Source = RssFeeds.Channel.Title
		item.SourceURL = RssFeeds.Channel.Link.Href
		item.Link = RssFeeds.Channel.Item[i].Link
		item.PublishDate, err = parseDate(RssFeeds.Channel.Item[i].PubDate)
		if err != nil {
			fmt.Print(err)
		}
		item.Description = RssFeeds.Channel.Item[i].Description
		items = append(items, item)
	}
	return items
}

var dateFormats = []string{
	time.RFC822,  // RSS
	time.RFC822Z, // RSS
	time.RFC3339, // Atom
	time.UnixDate,
	time.RubyDate,
	time.RFC850,
	time.RFC1123Z,
	time.RFC1123,
	time.ANSIC,
	"Mon, January 2 2006 15:04:05 -0700",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, Jan 2 2006 15:04:05 -0700",
	"Mon Jan 2 15:04 2006",
	"Mon Jan 02, 2006 3:04 pm",
	"Mon Jan 02 2006 15:04:05 -0700",
	"Mon Jan 02 2006 15:04:05 GMT-0700 (MST)",
	"Monday, January 2, 2006 03:04 PM",
	"Monday, January 2, 2006",
	"Monday, January 02, 2006",
	"Monday, 2 January 2006 15:04:05 -0700",
	"Monday, 2 Jan 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05",
	"Mon, 2 January 2006, 15:04 -0700",
	"Mon, 2 January 2006 15:04:05 -0700",
	"Mon, 2 January 2006",
	"Mon, 2 Jan 2006 3:04:05 PM -0700",
	"Mon, 2 Jan 2006 15:4:5 -0700 GMT",
	"Mon, 2, Jan 2006 15:4",
	"Mon, 2 Jan 2006, 15:04 -0700",
	"Mon, 2 Jan 2006 15:04 -0700",
	"Mon, 2 Jan 2006 15:04:05 UT",
	"Mon, 2 Jan 2006 15:04:05 -0700 MST",
	"Mon, 2 Jan 2006 15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05",
	"Mon, 2 Jan 2006 15:04",
	"Mon,2 Jan 2006",
	"Mon, 2 Jan 2006",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 2006-01-02 15:04",
	"Mon, 02 January 2006",
	"Mon, 02 Jan 2006 15 -0700",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04:05 Z",
	"Mon, 02 Jan 2006 15:04:05 UT",
	"Mon, 02 Jan 2006 15:04:05 MST-07:00",
	"Mon, 02 Jan 2006 15:04:05 MST -0700",
	"Mon, 02 Jan 2006 15:04:05 GMT-0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05 --0700",
	"Mon 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07",
	"Mon, 02 Jan 2006 15:04:05 00",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:4:5 Z",
	"Mon, 02 Jan 2006",
	"January 2, 2006 3:04 PM",
	"January 2, 2006, 3:04 p.m.",
	"January 2, 2006 15:04:05",
	"January 2, 2006 03:04 PM",
	"January 2, 2006",
	"January 02, 2006 15:04",
	"January 02, 2006 03:04 PM",
	"January 02, 2006",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2, 2006",
	"Jan 02 2006 03:04:05PM",
	"Jan 02, 2006",
	"6/1/2 15:04",
	"6-1-2 15:04",
	"2 January 2006 15:04:05 -0700",
	"2 January 2006",
	"2 Jan 2006 15:04:05 Z",
	"2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006",
	"2.1.2006 15:04:05",
	"2/1/2006",
	"2-1-2006",
	"2006 January 02",
	"2006-1-2T15:04:05Z",
	"2006-1-2 15:04:05",
	"2006-1-2",
	"2006-1-02T15:04:05Z",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00:00",
	"2006-01-02T15:04:05:-0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05 -0700",
	"2006-01-02T15:04:05:00",
	"2006-01-02T15:04:05",
	"2006-01-02 at 15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04",
	"2006-01-02 00:00:00.0 15:04:05.0 -0700",
	"2006/01/02",
	"2006-01-02",
	"15:04 02.01.2006 -0700",
	"1/2/2006 3:04:05 PM",
	"1/2/2006",
	"06/1/2 15:04",
	"06-1-2 15:04",
	"02 Monday, Jan 2006 15:04",
	"02 Jan 2006 15:04:05 UT",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05",
	"02 Jan 2006",
	"02.01.2006 15:04:05",
	"02/01/2006 15:04:05",
	"02.01.2006 15:04",
	"02/01/2006 - 15:04",
	"02.01.2006 -0700",
	"02/01/2006",
	"02-01-2006",
	"01/02/2006 3:04 PM",
	"01/02/2006 - 15:04",
	"01/02/2006",
	"01-02-2006",
}

// Named zone cannot be consistently loaded, so handle separately
var dateFormatsWithNamedZone = []string{
	"Mon, January 02, 2006, 15:04:05 MST",
	"Mon, January 02, 2006 15:04:05 MST",
	"Mon, Jan 2, 2006 15:04 MST",
	"Mon, Jan 2 2006 15:04 MST",
	"Mon, Jan 2, 2006 15:04:05 MST",
	"Mon Jan 2 15:04:05 2006 MST",
	"Mon, Jan 02,2006 15:04:05 MST",
	"Monday, January 2, 2006 15:04:05 MST",
	"Monday, 2 January 2006 15:04:05 MST",
	"Monday, 2 Jan 2006 15:04:05 MST",
	"Monday, 02 January 2006 15:04:05 MST",
	"Mon, 2 January 2006 15:04 MST",
	"Mon, 2 January 2006, 15:04:05 MST",
	"Mon, 2 January 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:4:5 MST",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 2006 15:04:05MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon 2 Jan 2006 15:04:05 MST",
	"mon,2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 15:04:05 MST",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon,02 January 2006 14:04:05 MST",
	"Mon, 02 Jan 2006 3:04:05 PM MST",
	"Mon,02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006, 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon , 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 06 15:04:05 MST",
	"January 2, 2006 15:04:05 MST",
	"January 02, 2006 15:04:05 MST",
	"Jan 2, 2006 3:04:05 PM MST",
	"Jan 2, 2006 15:04:05 MST",
	"2 January 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 MST",
	"2006-01-02 15:04:05 MST",
	"1/2/2006 3:04:05 PM MST",
	"1/2/2006 15:04:05 MST",
	"02 Jan 2006 15:04 MST",
	"02 Jan 2006 15:04:05 MST",
	"02/01/2006 15:04 MST",
	"02-01-2006 15:04:05 MST",
	"01/02/2006 15:04:05 MST",
}

// ParseDate parses a given date string using a large
// list of commonly found feed date formats.
func parseDate(ds string) (t time.Time, err error) {
	d := strings.TrimSpace(ds)
	if d == "" {
		return t, fmt.Errorf("date string is empty")
	}
	for _, f := range dateFormats {
		if t, err = time.Parse(f, d); err == nil {
			return
		}
	}
	for _, f := range dateFormatsWithNamedZone {
		t, err = time.Parse(f, d)
		if err != nil {
			continue
		}

		// This is a format match! Now try to load the timezone name
		loc, err := time.LoadLocation(t.Location().String())
		if err != nil {
			// We couldn't load the TZ name. Just use UTC instead...
			return t, nil
		}

		if t, err = time.ParseInLocation(f, ds, loc); err == nil {
			return t, nil
		}
		// This should not be reachable
	}

	err = fmt.Errorf("failed to parse date: %s", ds)
	return
}
