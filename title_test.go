package boondoggle

import (
	"bytes"
	"testing"
)

func TestExtractTitle(t *testing.T) {
	// Create mock articles
	var article Article

	// Nothing should do nothing
	article = Article{Raw: []byte(``)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "" {
		t.Errorf("Unexpected title: %s", article.Title)
	}

	// TODO table tests are funky looking with multiline strings
	// TODO ExtractTitle will add a newline - it shouldn'y
	md1 := `#Title
More Text 
`
	article = Article{Raw: []byte(md1)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "Title" {
		t.Errorf("Unexpected title: Title != %s", article.Title)
	}
	remainder := `More Text 
`
	if !bytes.Equal(article.Raw, []byte(remainder)) {
		t.Errorf("Unexpected raw buffer: `More Text ` != %s", article.Raw)
	}

	md2 := `# Title of Article 

Body text
`
	article = Article{Raw: []byte(md2)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "Title of Article" {
		t.Errorf("Unexpected title: Title of Article != %s", article.Title)
	}

	md3 := `❤ 
=====
 More
`
	article = Article{Raw: []byte(md3)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "❤" {
		t.Errorf("Unexpected title: ❤ != %s", article.Title)
	}
	remainder = ` More
` // TODO Fix newlines
	if !bytes.Equal(article.Raw, []byte(remainder)) {
		t.Errorf("Unexpected raw buffer: ` More` != %s", article.Raw)
	}

	md4 := `One line
`
	article = Article{Raw: []byte(md4)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "" {
		t.Errorf("Unexpected title: %s", article.Title)
	}

	md5 := `Just Text
Nothing more`
	article = Article{Raw: []byte(md5)}
	if err := ExtractTitle(&article); err != nil {
		t.Fatalf("ExtractTitle should not error: %s", err)
	}
	if article.Title != "" {
		t.Errorf("Unexpected title: %s", article.Title)
	}
	if !bytes.Equal(article.Raw, []byte(md5)) {
		t.Errorf("Unexpected raw buffer: %s != %s", md5, article.Raw)
	}
}
