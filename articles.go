package boondoggle

import "sort"

// Articles is a slice of articles
type Articles []Article

// Len returns the number of articles
func (a Articles) Len() int {
	return len(a)
}

// Swap is used to sort articles
func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ByTitle will sort the articles by slug
type ByTitle struct {
	Articles
}

// ByDate will sort articles by timestamp
type ByDate struct {
	Articles
}

// Implement the sort.Interface for sorting
var _ sort.Interface = ByDate{}

func (a ByDate) Less(i, j int) bool {
	x, y := a.Articles[i], a.Articles[j]
	if x.Date.Unix() == y.Date.Unix() {
		// Sort alphabetically
		return x.Title < y.Title
	}
	// Most recent articles should be first
	return x.Date.Unix() > y.Date.Unix()
}

// SortMostRecentArticlesFirst will sort the articles by Date
func (a Articles) SortMostRecentArticlesFirst() {
	sort.Sort(ByDate{a})
}
