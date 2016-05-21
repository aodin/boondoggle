package boondoggle

import "testing"

// Test the sorting of articles
func TestArticles(t *testing.T) {
	a := Article{
		Title: "A",
		Date:  MustCreate("2013-10-30"),
	}
	b := Article{
		Title: "b",
		Date:  MustCreate("2013-10-30"),
	}
	x := Article{
		Title: "X",
		Date:  MustCreate("2013-01-30"),
	}

	articles := Articles{x, b, a}
	articles.SortByDate()

	if articles[0].Title != "A" {
		t.Error("Unexpected sort order for articles:", articles)
	}
	if articles[1].Title != "b" {
		t.Error("Unexpected sort order for articles:", articles)
	}
	if articles[2].Title != "X" {
		t.Error("Unexpected sort order for articles:", articles)
	}
}
