package boondoggle

import (
	"bytes"
	"testing"
)

var fmExample = `---
title: I am an article
tags: sql, DATA SCIENCE , python
---
Stuff
`

func TestExtractFrontMatter(t *testing.T) {
	// Create mock articles
	var article Article

	// Nothing should do nothing
	article = Article{Raw: []byte(``)}
	if err := ExtractFrontMatter(&article); err != nil {
		t.Fatalf("ExtractFrontMatter should not error: %s", err)
	}
	if article.Title != "" {
		t.Errorf("Unexpected title: %s", article.Title)
	}

	// Empty blocks should not error
	article = Article{Raw: []byte(`---
---
Yo`)}
	if err := ExtractFrontMatter(&article); err != nil {
		t.Fatalf("ExtractFrontMatter should not error: %s", err)
	}

	// Real example
	article = Article{Raw: []byte(fmExample)}
	if err := ExtractFrontMatter(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "I am an article" {
		t.Errorf("Unexpected title: I am an article != %s", article.Title)
	}
	if len(article.Tags) != 3 {
		t.Fatalf("Unexpected number of tags: 3 != %d", len(article.Tags))
	}
	remainder := `Stuff
`
	if !bytes.Equal(article.Raw, []byte(remainder)) {
		t.Errorf("Unexpected raw buffer: `Stuff` != %s", article.Raw)
	}
}
