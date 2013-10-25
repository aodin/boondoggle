package boondoggle

import (
	"testing"
)

var testDir = `./testdata/`

var example = `
Boondoggle
==========

Subtitle
--------

### Example Title

    Some(code) = this

#### List

* Just
* A List
* Of Things

More on GitHub, [Boondoggle](www.github.com/aodin/boondoggle)
`

func TestLoadArticles(t *testing.T) {
	articles, err := LoadArticles(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 2 {
		t.Fatalf("Unexpected length of articles: %d", len(articles))
	}
	article := articles[1]
	if len(article.Raw) != 107 {
		t.Errorf("Unexpected raw content length: %s", len(article.Raw))
	}
	if len(article.Body) != 162 {
		t.Errorf("Unexpected body length: %s", len(article.Body))
	}
	if article.Title != "Second Post" {
		t.Errorf("Unexpected article title: %s", article.Title)
	}
}

// Test the sorting of articles
func TestArticles(t *testing.T) {
	a := &Article{
		Title: "A",
		Date:  MustCreate("2013-10-30"),
	}
	b := &Article{
		Title: "B",
		Date:  MustCreate("2013-10-30"),
	}
	x := &Article{
		Title: "X",
		Date:  MustCreate("2013-01-30"),
	}

	articles := Articles{x, b, a}
	articles.Sort()

	if articles[0] != a {
		t.Error("Unexpected sort order for articles:", articles)
	}
	if articles[1] != b {
		t.Error("Unexpected sort order for articles:", articles)
	}
	if articles[2] != x {
		t.Error("Unexpected sort order for articles:", articles)
	}
}
