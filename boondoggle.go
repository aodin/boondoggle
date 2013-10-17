package boondoggle

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/russross/blackfriday"
)

var slugifyClean = regexp.MustCompile(`[^\w\s-]`)
var slugifySpace = regexp.MustCompile(`[-\s]+`)

type Article struct {
	Title string
	Body  template.HTML
	Raw   []byte // The raw markdown string
}

func (article *Article) String() string {
	return article.Title
}

func Slugify(input string) string {
	// Remove anything that isn't a digit, word character or dash
	lowered := string(slugifyClean.ReplaceAll([]byte(input), []byte("")))
	lowered = strings.TrimSpace(strings.ToLower(lowered))
	return string(slugifySpace.ReplaceAll([]byte(lowered), []byte("-")))
}

func UnSnakeCase(input string) string {
	words := strings.Split(input, "_")
	titles := make([]string, len(words))
	// TODO Inefficient method for capitalization, look at unicode.ToTitle
	for index, word := range words {
		titles[index] = strings.Title(word)
	}
	return strings.Join(titles, " ")
}

func GetArticles(path string) (articles []*Article, err error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		// Easy way to get the fullpath
		name := entry.Name()
		extension := strings.ToLower(filepath.Ext(name))

		// TODO Allow a filter by file extension
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

		// TODO We'll cheat for now to get the filename without the extension,
		// since we know it ends in .md
		article := &Article{
			Title: UnSnakeCase(name[:len(name) - 3]),
			Raw: content,
			Body: template.HTML(blackfriday.MarkdownBasic(content)),
		}
		articles = append(articles, article)
	}
	return
}
