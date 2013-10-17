package boondoggle

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/russross/blackfriday"
)

type Article struct {
	Title string
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
