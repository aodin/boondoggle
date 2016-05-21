/*
Boondoggle is a static site generator written in Go.
*/
package boondoggle

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Boondoggle builds .HTML files from a directory of markdown files.
type Boondoggle struct {
	articles map[string]Article   // By title - including path!
	tags     map[string][]Article // TODO pointer?

	listTemplate    *template.Template
	listCache       []byte
	articleTemplate *template.Template
	ordering        []Article
	attrs           map[string]interface{}
	logger          RequestLogger
}

// TODO AllTags returns tags in alphabetical order

func (b *Boondoggle) ArticleTemplate(tmpl *template.Template) *Boondoggle {
	b.articleTemplate = tmpl
	return b
}

func (b *Boondoggle) ListTemplate(tmpl *template.Template) *Boondoggle {
	b.listTemplate = tmpl
	return b
}

func (b *Boondoggle) Attr(key string, value interface{}) *Boondoggle {
	b.attrs[key] = value
	return b
}

func (b *Boondoggle) LoadFrom(path string) error {
	articles, err := LoadArticles(path)
	if err != nil {
		return err
	}
	b.ordering = articles
	Articles(articles).SortByDate()
	for _, article := range articles {
		// TODO What to do about duplicate slugs?
		b.articles[article.Slug] = article
	}
	log.Printf("Loaded %d Articles\n", len(articles))
	return nil
}

// Route to the requested article, if it exists
func (b *Boondoggle) Route(w http.ResponseWriter, r *http.Request) {
	// Log the request
	b.logger.Log(r)

	// We assume the last part of the request URL is the article slug
	path := strings.Split(r.URL.Path, "/")
	slug := path[len(path)-1]
	article, exists := b.articles[slug]
	if exists {
		b.Article(w, article)
		return
	}

	// List the articles if it was an empty path
	if slug == "" {
		b.List(w)
		return
	}

	// Redirect to the List
	http.Redirect(w, r, path[0], 302)
}

func (b *Boondoggle) Article(w http.ResponseWriter, article Article) {
	if len(article.Cache) == 0 {
		b.attrs["Article"] = article
		buffer := &bytes.Buffer{}
		b.articleTemplate.Execute(buffer, b.attrs)
		article.Cache = buffer.Bytes()
	}
	w.Write(article.Cache)
}

// List the available articles
func (b *Boondoggle) List(w http.ResponseWriter) {
	if len(b.listCache) == 0 {
		b.attrs["Articles"] = b.ordering
		buffer := &bytes.Buffer{}
		b.listTemplate.Execute(buffer, b.attrs)
		b.listCache = buffer.Bytes()
	}
	w.Write(b.listCache)
}

// TODO Or this could return the handler
func (b *Boondoggle) Handler(w http.ResponseWriter, r *http.Request) {
	b.Route(w, r)
}

var listTmpl = `<!DOCTYPE html>
<html>
  <head>
    <title>Articles</title>
  </head>
  <body>
    <h1>Articles</h1>
    <ul>
    {{ range $article := .Articles }}
      <li><a href="./{{ $article.Slug }}">{{ $article.Title }}</a></li>{{ end }}
    </ul>
  </body>
</html>
`

var articleTmpl = `<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Article.Title }}</title>
  </head>
  <body>
  	<h1>{{ .Article.Title }}</h1>
    {{ .Article.Body }}
  </body>
</html>
`

// Create an empty boondoggle
func Create() *Boondoggle {
	return &Boondoggle{
		articles:        make(map[string]Article),
		listTemplate:    template.Must(template.New("list").Parse(listTmpl)),
		articleTemplate: template.Must(template.New("article").Parse(articleTmpl)),
		listCache:       make([]byte, 0),
		attrs:           make(map[string]interface{}),
		logger:          defaultLogger,
	}
}

// Create a boondoggle by loading articles from the given directory
func CreateFrom(path string) (*Boondoggle, error) {
	// Load the articles from the directory
	b := Create()
	err := b.LoadFrom(path)
	if err != nil {
		return b, err
	}
	return b, nil
}
