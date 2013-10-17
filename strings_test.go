package boondoggle

import (
	"testing"
)

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

	input = "  LOTS of SPACES  "
	output = Slugify(input)
	expected = "lots-of-spaces"
	if output != expected {
		t.Errorf("Unexpected Slugify() output: %s != %s", output, expected)
	}
}
