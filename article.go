package boondoggle

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Article is a single article
type Article struct {
	// Content
	Title    string
	Slug     string
	Date     Timestamp // TODO time.Time
	Subtitle string
	Byline   string
	HTML     template.HTML

	// Meta
	WordCount       uint64
	TableOfContents TableOfContents
	LinesOfCode     uint64
	Tags            []string

	// TODO need methods to create buffers/scanners and reset raw
	Body  template.HTML // Deprecated
	Raw   []byte        // The entire raw file - TODO un-exported
	Cache []byte        // The executed template - TODO delete
}

func (article *Article) String() string {
	return article.Title
}

func ParseArticle(content []byte) Article {
	// Remove the title and an optional date from the content
	index := 0
	last := 0
	length := len(content)

	// TODO Named return type won't work?
	article := Article{Raw: content}

	var headers [3][]byte
	for header, _ := range headers {
		for index < length {
			// TODO robust newline handling
			if content[index] == '\n' {
				index += 1
				break
			}
			index += 1
		}
		headers[header] = content[last : index-1]
		last = index
	}
	article.Title = string(headers[0])
	article.Subtitle = string(headers[1])
	article.Byline = string(headers[2])
	article.Body = template.HTML(blackfriday.MarkdownCommon(content[last:]))
	return article
}

func LoadArticles(path string) (articles []Article, err error) {
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
		content, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return nil, err
		}

		article := ParseArticle(content)

		// Attempt to parse a date from the filename
		// TODO We'll cheat for now to get the filename without the extension,
		// since we know it ends in .md
		filename := name[:len(name)-3]

		date, slug := SplitFilename(filename)
		// If not date was recovered, convert the whole filename to a slug
		if date == "" {
			slug = filename
		} else {
			// Attempt to parse the date and add it to the article
			timestamp, err := CreateTimestamp(date)
			if err == nil {
				article.Date = timestamp
			}
		}
		article.Slug = slug
		articles = append(articles, article)
	}
	return
}
