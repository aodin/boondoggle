package boondoggle

import (
	"bytes"
	"testing"
)

func TestParseTags(t *testing.T) {
	// Create mock articles
	var article Article

	// Nothing should do nothing
	article = Article{Raw: []byte(``)}
	if err := ParseTags(&article); err != nil {
		t.Fatalf("ParseTags should not error: %s", err)
	}
	if len(article.Tags) != 0 {
		t.Errorf("Unexpected number of tags: %d", len(article.Tags))
	}

	// TODO table tests are funky looking with multiline strings
	md1 := `#Title

[//]: # (golang, SQL ,data science)  

More Text `
	article = Article{Raw: []byte(md1)}
	if err := ParseTags(&article); err != nil {
		t.Fatalf("ParseTags should not error: %s", err)
	}
	if len(article.Tags) != 3 {
		t.Fatalf("Unexpected number of tags: 3 != %d", len(article.Tags))
	}
	if article.Tags[1] != "sql" {
		t.Errorf("Unexpected tag: sql != %s", article.Tags[1])
	}
	if article.Tags[2] != "data-science" {
		t.Errorf("Unexpected tag: data-science != %s", article.Tags[2])
	}
	// The raw bytes should not be modified
	if !bytes.Equal(article.Raw, []byte(md1)) {
		t.Errorf("Unexpected raw buffer: %s != %s", md1, article.Raw)
	}
}
