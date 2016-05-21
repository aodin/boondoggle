package boondoggle

import (
	"testing"
	"time"
)

func TestParseFilename(t *testing.T) {
	// Create mock articles
	var article Article

	article = Article{}
	if err := ParseFilename(&article); err != nil {
		t.Fatalf("ParseFilename should not error: %s", err)
	}
	if !article.Date.IsZero() {
		t.Errorf("Unexpected date: %s", article.Date)
	}

	article = Article{Filename: "A Word.md"}
	if err := ParseFilename(&article); err != nil {
		t.Fatalf("ParseFilename should not error: %s", err)
	}
	if article.Title != "A Word" {
		t.Errorf("Unexpected title: A Word != %s", article.Title)
	}
	if article.Slug != "a-word" {
		t.Errorf("Unexpected slug: a-word != %s", article.Slug)
	}
	if !article.Date.IsZero() {
		t.Errorf("Unexpected date: %s", article.Date)
	}

	article = Article{Filename: "Some Post About _.md"}
	if err := ParseFilename(&article); err != nil {
		t.Fatalf("ParseFilename should not error: %s", err)
	}
	if article.Title != "Some Post About _" {
		t.Errorf("Unexpected title: Some Post About _ != %s", article.Title)
	}
	if article.Slug != "some-post-about" {
		t.Errorf("Unexpected slug: some-post-about != %s", article.Slug)
	}
	if !article.Date.IsZero() {
		t.Errorf("Unexpected date: %s", article.Date)
	}

	article = Article{Filename: "2016-03-01_post"}
	if err := ParseFilename(&article); err != nil {
		t.Fatalf("ParseFilename should not error: %s", err)
	}
	if article.Title != "post" {
		t.Errorf("Unexpected title: post != %s", article.Title)
	}
	if article.Slug != "post" {
		t.Errorf("Unexpected slug: post != %s", article.Slug)
	}
	expected := time.Date(2016, 3, 1, 0, 0, 0, 0, time.UTC)
	if !article.Date.Equal(expected) {
		t.Errorf("Unexpected date: %s != %s", article.Date, expected)
	}
}
