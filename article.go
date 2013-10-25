package boondoggle

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Article struct {
	Title string
	Slug  string
	Date  Timestamp
	Body  template.HTML
	Raw   []byte // The raw markdown string
}

func (article *Article) String() string {
	return article.Title
}

func LoadArticles(path string) (articles []*Article, err error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		// TODO Is there an easier way to get the fullpath?
		// TODO Walk the entire directory structure?
		name := entry.Name()
		extension := strings.ToLower(filepath.Ext(name))

		// TODO Allow a filter by custom file extension
		if extension != ".md" {
			continue
		}
		fullpath := filepath.Join(path, name)
		file, err := os.Open(fullpath)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		// Attempt to parse a date from the filename
		// TODO We'll cheat for now to get the filename without the extension,
		// since we know it ends in .md
		filename := name[:len(name)-3]

		date, title := SplitFilename(filename)
		// If not date was recovered, convert the whole filename to a title
		if date == "" {
			title = filename
		}

		article := &Article{
			Title: UnSnakeCase(title),
			Raw:   content,
			Body:  template.HTML(blackfriday.MarkdownCommon(content)),
		}

		// TODO guh, stupid logical flow
		if date != "" {
			// Attempt to parse the date and add it to the article
			timestamp, err := CreateTimestamp(date)
			if err == nil {
				article.Date = timestamp
			}
		}

		articles = append(articles, article)
	}
	return
}

type Articles []*Article

// Implement the sort.Interface for sorting
func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByTitle struct {
	Articles
}

type ByDate struct {
	Articles
}

func (a ByDate) Less(i, j int) bool {
	x, y := a.Articles[i], a.Articles[j]
	if x.Date.Unix() == y.Date.Unix() {
		// Sort alphabetically
		return x.Title < y.Title
	}
	// Most recent articles should be first
	return x.Date.Unix() > y.Date.Unix()
}

func (a Articles) Sort() {
	sort.Sort(ByDate{a})
}
