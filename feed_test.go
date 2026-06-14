package boondoggle

import (
	"encoding/xml"
	"html/template"
	"strings"
	"testing"
	"time"
)

// sampleBoondoggle builds a Boondoggle with two articles for feed testing.
func sampleBoondoggle() *Boondoggle {
	bd := New()
	bd.Articles = Articles{
		{
			Title:   "Second Post",
			Slug:    "second-post",
			Date:    time.Date(2024, 2, 1, 9, 0, 0, 0, time.UTC),
			HTML:    template.HTML("<p>The <em>second</em> post body.</p>"),
			Preview: template.HTML("<p>The second post body.</p>"),
			Links:   bd.Links,
		},
		{
			Title: "First Post",
			Slug:  "first-post",
			Date:  time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			HTML:  template.HTML("<p>The first post body & more.</p>"),
			Links: bd.Links,
		},
	}
	// Feeds expect the most recent article first.
	bd.Articles.SortMostRecentArticlesFirst()
	return bd
}

func sampleFeed() Feed {
	return Feed{
		Title:       "Example Blog",
		Link:        "https://example.com",
		Description: "An example blog",
		Author:      "Jane Doe",
		Email:       "jane@example.com",
	}
}

func TestFeedRSS(t *testing.T) {
	out, err := sampleFeed().RSS(sampleBoondoggle())
	if err != nil {
		t.Fatalf("RSS returned an error: %s", err)
	}

	// The output must be well-formed XML that parses back into an RSS document.
	var feed rssFeed
	if err := xml.Unmarshal(out, &feed); err != nil {
		t.Fatalf("RSS output is not valid XML: %s\n%s", err, out)
	}

	if feed.Version != "2.0" {
		t.Errorf("Unexpected RSS version: %q", feed.Version)
	}
	if feed.Channel.Title != "Example Blog" {
		t.Errorf("Unexpected channel title: %q", feed.Channel.Title)
	}
	if feed.Channel.Link != "https://example.com" {
		t.Errorf("Unexpected channel link: %q", feed.Channel.Link)
	}
	if len(feed.Channel.Items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(feed.Channel.Items))
	}

	// Most recent article should be first.
	first := feed.Channel.Items[0]
	if first.Title != "Second Post" {
		t.Errorf("Expected most recent item first, got %q", first.Title)
	}
	if first.Link != "https://example.com/articles/second-post" {
		t.Errorf("Unexpected item link: %q", first.Link)
	}
	if first.GUID != first.Link {
		t.Errorf("Expected guid to equal link, got %q", first.GUID)
	}
	if want := "Thu, 01 Feb 2024 09:00:00 +0000"; first.PubDate != want {
		t.Errorf("Unexpected pubDate: %q (want %q)", first.PubDate, want)
	}
	if first.Author != "jane@example.com (Jane Doe)" {
		t.Errorf("Unexpected author: %q", first.Author)
	}

	// HTML content survives a round-trip through escaped XML.
	second := feed.Channel.Items[1]
	if second.Description != "<p>The first post body & more.</p>" {
		t.Errorf("Unexpected description: %q", second.Description)
	}

	// The raw bytes should carry the XML declaration and escape the ampersand.
	if !strings.HasPrefix(string(out), "<?xml") {
		t.Errorf("Missing XML declaration:\n%s", out)
	}
	if !strings.Contains(string(out), "&amp;") {
		t.Errorf("Expected ampersand to be escaped:\n%s", out)
	}
}

func TestFeedAtom(t *testing.T) {
	out, err := sampleFeed().Atom(sampleBoondoggle())
	if err != nil {
		t.Fatalf("Atom returned an error: %s", err)
	}

	var feed atomFeed
	if err := xml.Unmarshal(out, &feed); err != nil {
		t.Fatalf("Atom output is not valid XML: %s\n%s", err, out)
	}

	if feed.Title != "Example Blog" {
		t.Errorf("Unexpected feed title: %q", feed.Title)
	}
	if feed.ID != "https://example.com" {
		t.Errorf("Unexpected feed id: %q", feed.ID)
	}
	// Updated should reflect the most recent article.
	if want := "2024-02-01T09:00:00Z"; feed.Updated != want {
		t.Errorf("Unexpected feed updated: %q (want %q)", feed.Updated, want)
	}
	if feed.Author == nil || feed.Author.Name != "Jane Doe" || feed.Author.Email != "jane@example.com" {
		t.Errorf("Unexpected feed author: %+v", feed.Author)
	}
	if len(feed.Entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(feed.Entries))
	}

	first := feed.Entries[0]
	if first.Title != "Second Post" {
		t.Errorf("Expected most recent entry first, got %q", first.Title)
	}
	if first.ID != "https://example.com/articles/second-post" {
		t.Errorf("Unexpected entry id: %q", first.ID)
	}
	if len(first.Links) != 1 || first.Links[0].Href != first.ID {
		t.Errorf("Unexpected entry links: %+v", first.Links)
	}
	if first.Published != "2024-02-01T09:00:00Z" {
		t.Errorf("Unexpected entry published: %q", first.Published)
	}
	if first.Content == nil || first.Content.Type != "html" {
		t.Fatalf("Unexpected entry content: %+v", first.Content)
	}
	if first.Content.Body != "<p>The <em>second</em> post body.</p>" {
		t.Errorf("Unexpected entry content body: %q", first.Content.Body)
	}
	if first.Summary == nil || first.Summary.Body != "<p>The second post body.</p>" {
		t.Errorf("Unexpected entry summary: %+v", first.Summary)
	}

	// The second entry has no preview, so it should omit the summary element.
	if feed.Entries[1].Summary != nil {
		t.Errorf("Expected no summary for article without preview")
	}

	if !strings.Contains(string(out), `xmlns="http://www.w3.org/2005/Atom"`) {
		t.Errorf("Missing Atom namespace:\n%s", out)
	}
}

func TestFeedLimit(t *testing.T) {
	feed := sampleFeed()
	feed.Limit = 1

	out, err := feed.RSS(sampleBoondoggle())
	if err != nil {
		t.Fatalf("RSS returned an error: %s", err)
	}

	var parsed rssFeed
	if err := xml.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("RSS output is not valid XML: %s", err)
	}
	if len(parsed.Channel.Items) != 1 {
		t.Errorf("Expected limit to cap items at 1, got %d", len(parsed.Channel.Items))
	}
}

// The renderers must satisfy the FeedTransformer signature so they can be used
// interchangeably.
func TestFeedTransformerSignature(t *testing.T) {
	transformers := map[string]FeedTransformer{
		"rss":  sampleFeed().RSS,
		"atom": sampleFeed().Atom,
	}
	for name, transform := range transformers {
		out, err := transform(sampleBoondoggle())
		if err != nil {
			t.Fatalf("%s transformer returned an error: %s", name, err)
		}
		if len(out) == 0 {
			t.Errorf("%s transformer produced no output", name)
		}
	}
}
