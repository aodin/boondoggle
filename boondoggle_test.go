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

func TestGetArticles(t *testing.T) {
	articles, err := GetArticles(testDir)
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

func TestUnSnakeCase(t *testing.T) {
	input := "snake_case_title"
	output := UnSnakeCase(input)
	expected := "Snake Case Title"
	if output != expected {
		t.Errorf("Unexpected UnSnakeCase() output: %s != %s", output, expected)
	}
}

func TestSlugify(t *testing.T) {
	input := "Article Title"
	output := Slugify(input)
	expected := "article-title"
	if output != expected {
		t.Errorf("Unexpected Slugify() output: %s != %s", output, expected)
	}
}