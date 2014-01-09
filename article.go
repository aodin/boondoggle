package boondoggle

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type Article struct {
	Slug     string
	Date     Timestamp
	Title    string
	Subtitle string
	Byline   string
	Body     template.HTML
	Raw      []byte // The entire raw file
	Cache    []byte // The executed template
}

func (article *Article) String() string {
	return article.Title
}

func ParseArticle(content []byte) *Article {
	// Remove the title and an optional date from the content
	index := 0
	last := 0
	length := len(content)

	// TODO Named return type won't work?
	article := &Article{Raw: content}

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
