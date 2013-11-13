package boondoggle

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Boondoggle struct {
	articles        map[string]*Article
	listTemplate    *template.Template
	articleTemplate *template.Template
	ordering        []*Article
	attrs           map[string]interface{}
}

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
	Articles(articles).Sort()
	for _, article := range articles {
		article.Slug = Slugify(article.Title)
		// TODO What to do about duplicate slugs?
		b.articles[article.Slug] = article
	}
	log.Printf("Loaded %d Articles\n", len(articles))
	return nil
}

// Route to the requested article, if it exists
func (b *Boondoggle) Route(w http.ResponseWriter, r *http.Request) {
	// Log the request
	Log(r)

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
	// TODO use closure around a copy of the attrs? Cache the result?
	b.attrs["Article"] = article
	b.articleTemplate.Execute(w, b.attrs)
}

// List the available articles
func (b *Boondoggle) List(w http.ResponseWriter) {
	// TODO There should be a cached list of articles
	b.attrs["Articles"] = b.ordering
	b.listTemplate.Execute(w, b.attrs)
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
    {{ .Article.Body }}
  </body>
</html>
`

// Create an empty boondoggle
func Create() *Boondoggle {
	return &Boondoggle{
		articles:        make(map[string]*Article),
		listTemplate:    template.Must(template.New("list").Parse(listTmpl)),
		articleTemplate: template.Must(template.New("article").Parse(articleTmpl)),
		attrs:           make(map[string]interface{}),
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

// TODO Logging should be an interface
func Log(r *http.Request) {
	// By default, a timestamp will be written by the logger
	log.Printf(`"%s %s" %s "%s"`, r.Method, r.URL, strings.SplitN(r.RemoteAddr, ":", 2)[0], r.Header.Get("User-Agent"))
}
