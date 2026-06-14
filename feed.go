package boondoggle

import (
	"encoding/xml"
	"strings"
	"time"
)

// FeedTransformer produces a serialized syndication feed (such as RSS or Atom
// XML) from a parsed Boondoggle. Unlike a Transformer, which operates on a
// single Article, FeedTransformer aggregates every article into a single document.
type FeedTransformer func(*Boondoggle) ([]byte, error)

// Feed holds site-level metadata required to render a syndication feed.
// Per-article values (title, link, date, content) are taken from the articles.
type Feed struct {
	Title       string `toml:"title"`
	Link        string `toml:"url"` // Absolute base URL of the site, e.g. https://example.com
	Description string `toml:"description"`
	Author      string `toml:"author"`
	Email       string `toml:"email"`
	Limit       int    `toml:"items"` // Set 0 to include all articles
}

// RSS renders the Boondoggle's articles as an RSS 2.0 feed.
func (f Feed) RSS(bd *Boondoggle) ([]byte, error) {
	channel := rssChannel{
		Title:       f.Title,
		Link:        f.absURL(""),
		Description: f.Description,
	}
	if t := f.updated(bd); !t.IsZero() {
		channel.LastBuildDate = t.Format(time.RFC1123Z)
	}

	for _, article := range f.articles(bd) {
		link := f.absURL(article.URL())
		item := rssItem{
			Title:       article.Title,
			Link:        link,
			GUID:        link,
			Author:      f.rssAuthor(),
			Description: string(article.HTML),
		}
		if !article.Date.IsZero() {
			item.PubDate = article.Date.Format(time.RFC1123Z)
		}
		channel.Items = append(channel.Items, item)
	}

	return marshalFeed(rssFeed{Version: "2.0", Channel: channel})
}

// Atom renders the Boondoggle's articles as an Atom 1.0 feed.
func (f Feed) Atom(bd *Boondoggle) ([]byte, error) {
	self := f.absURL("")
	feed := atomFeed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Title:   f.Title,
		ID:      self,
		Updated: f.updated(bd).Format(time.RFC3339),
		Links: []atomLink{
			{Href: self, Rel: "alternate"},
		},
	}
	if f.Description != "" {
		feed.Subtitle = f.Description
	}
	if f.Author != "" || f.Email != "" {
		feed.Author = &atomAuthor{Name: f.Author, Email: f.Email}
	}

	for _, article := range f.articles(bd) {
		link := f.absURL(article.URL())
		// Atom requires a valid updated timestamp; fall back to the build time
		// for articles that have no date.
		updated := article.Date
		if updated.IsZero() {
			updated = bd.BuildTime
		}
		entry := atomEntry{
			Title:   article.Title,
			ID:      link,
			Updated: updated.Format(time.RFC3339),
			Links:   []atomLink{{Href: link, Rel: "alternate"}},
			Content: &atomText{Type: "html", Body: string(article.HTML)},
		}
		if !article.Date.IsZero() {
			entry.Published = article.Date.Format(time.RFC3339)
		}
		if article.Preview != "" {
			entry.Summary = &atomText{Type: "html", Body: string(article.Preview)}
		}
		feed.Entries = append(feed.Entries, entry)
	}

	return marshalFeed(feed)
}

// Compile-time checks that the feed renderers satisfy FeedTransformer.
var (
	_ FeedTransformer = Feed{}.RSS
	_ FeedTransformer = Feed{}.Atom
)

// articles returns the (optionally limited) articles to include in the feed.
func (f Feed) articles(bd *Boondoggle) Articles {
	articles := bd.Articles
	if f.Limit > 0 && len(articles) > f.Limit {
		articles = articles[:f.Limit]
	}
	return articles
}

// updated returns the timestamp the feed was last updated, preferring the most
// recent article's date and falling back to the Boondoggle's build time.
func (f Feed) updated(bd *Boondoggle) time.Time {
	// Articles are sorted most recent first, so the first article is newest.
	if articles := f.articles(bd); len(articles) > 0 && !articles[0].Date.IsZero() {
		return articles[0].Date
	}
	return bd.BuildTime
}

// absURL joins the feed's base URL with a (possibly relative) path. Absolute
// paths are returned unchanged.
func (f Feed) absURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if path == "" {
		return f.Link
	}
	return strings.TrimRight(f.Link, "/") + "/" + strings.TrimLeft(path, "/")
}

// rssAuthor formats the author per the RSS spec, which expects an email address
// optionally followed by the name in parentheses. It returns an empty string
// when no email is configured.
func (f Feed) rssAuthor() string {
	if f.Email == "" {
		return ""
	}
	if f.Author != "" {
		return f.Email + " (" + f.Author + ")"
	}
	return f.Email
}

// marshalFeed serializes v to indented XML prefixed with the XML declaration.
func marshalFeed(v interface{}) ([]byte, error) {
	body, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	out := append([]byte(xml.Header), body...)
	return append(out, '\n'), nil
}

// RSS 2.0 document structure.

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Author      string `xml:"author,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
	Description string `xml:"description"`
}

// Atom 1.0 document structure.

type atomFeed struct {
	XMLName  xml.Name    `xml:"feed"`
	XMLNS    string      `xml:"xmlns,attr"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle,omitempty"`
	ID       string      `xml:"id"`
	Updated  string      `xml:"updated"`
	Links    []atomLink  `xml:"link"`
	Author   *atomAuthor `xml:"author,omitempty"`
	Entries  []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type atomAuthor struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
}

// atomText is an Atom text construct (used for content and summary). The Type
// attribute is typically "text" or "html".
type atomText struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type atomEntry struct {
	Title     string     `xml:"title"`
	ID        string     `xml:"id"`
	Updated   string     `xml:"updated"`
	Published string     `xml:"published,omitempty"`
	Links     []atomLink `xml:"link"`
	Summary   *atomText  `xml:"summary,omitempty"`
	Content   *atomText  `xml:"content,omitempty"`
}
