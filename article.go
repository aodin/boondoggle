package boondoggle

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	Title     string
	Date      string
	Timestamp time.Time
	Body      template.HTML
	Raw       []byte // The raw markdown string
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
			timestamp, err := ParseDate(date)
			if err == nil {
				article.Timestamp = timestamp
				article.Date = OutputDate(timestamp)
			}
		}

		articles = append(articles, article)
	}
	return
}
