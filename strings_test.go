package boondoggle

import (
	"testing"
	"time"
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

func TestSplitFilename(t *testing.T) {
	filename := `2013-10-22_hello_boondoggle`
	date, title := SplitFilename(filename)
	if date != "2013-10-22" {
		t.Errorf("Unexpected date part: %s != %s", date, "2013-10-22")
	}
	if title != "hello_boondoggle" {
		t.Errorf("Unexpected title part: %s != %s", title, "hello_boondoggle")
	}
}

func TestParseDate(t *testing.T) {
	input := "2013-10-22"
	output, err := ParseDate(input)
	if err != nil {
		t.Error("Error during date parsing:", err)
	}
	expect := time.Date(2013, time.October, 22, 0, 0, 0, 0, time.UTC)
	if output != expect {
		t.Errorf("Unexpected date: %s != %s", output.String(), expect.String())
	}
}

func TestOutputDate(t *testing.T) {
	input := time.Date(2013, time.October, 22, 0, 0, 0, 0, time.UTC)
	output := OutputDate(input)
	expect := `Tuesday, October 22, 2013`
	if output != expect {
		t.Errorf("Unexpected output: %s != %s", output, expect)
	}
}
