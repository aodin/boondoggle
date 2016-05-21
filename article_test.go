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
	if len(articles) != 3 {
		t.Fatalf("Unexpected length of articles: 3 != %d", len(articles))
	}
	article := articles[1]
	if len(article.Raw) != 107 {
		t.Errorf("Unexpected raw content length: %d", len(article.Raw))
	}
	if len(article.Body) != 140 {
		t.Errorf("Unexpected body length: %d", len(article.Body))
	}
	if article.Title != "Second Post" {
		t.Errorf("Unexpected article title: %s", article.Title)
	}
}
