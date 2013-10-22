package boondoggle

import (
	"html/template"
	"net/http"
	"strings"
)

// TODO Use multiple templates
// A basic template for the list
// TODO The template variables should be attached to the boondoggle, this
// will also allow them to be overwritten
var articleTemplate = template.Must(template.New("article").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Article.Title }}</title>
  </head>
  <body>
  	{{ .Article.Body }}
  </body>
 </html>
`))

var listTemplate = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>Articles</title>
  </head>
  <body>
  	<h1>Articles</h1>
  	<ul>
  	{{ range $slug, $article := .Articles }}
  		<li><a href="./{{ $slug }}">{{ $article.Title }}</a></li>{{ end }}
  	</ul>
  </body>
 </html>
`))

type Boondoggle struct {
	articles map[string]*Article
	// TODO Ordering
}

// Route to the requested article, if it exists
func (b *Boondoggle) Route(w http.ResponseWriter, r *http.Request) {
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
	http.NotFound(w, r)
}

func (b *Boondoggle) Article(w http.ResponseWriter, article *Article) {
	articleTemplate.Execute(w, map[string]interface{}{"Article": article})
}

// List the available articles
func (b *Boondoggle) List(w http.ResponseWriter) {
	// TODO There should be a cached list of articles
	listTemplate.Execute(w, map[string]interface{}{"Articles": b.articles})
}

// TODO Or this could return the handler
func (b *Boondoggle) Handler(w http.ResponseWriter, r *http.Request) {
	b.Route(w, r)
}

func CreateFrom(path string) (*Boondoggle, error) {
	// Load the articles from the directory
	b := &Boondoggle{
		articles: make(map[string]*Article),
	}
	articles, err := LoadArticles(path)
	if err != nil {
		return b, err
	}
	for _, article := range articles {
		slug := Slugify(article.Title)
		// TODO What to do about duplicate slugs?
		b.articles[slug] = article
	}
	return b, nil
}
