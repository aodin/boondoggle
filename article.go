package boondoggle

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// Article is a single article
type Article struct {
	// Content
	Title    string
	Slug     string
	Date     time.Time
	Subtitle string
	Byline   string
	HTML     template.HTML

	// Meta
	Preview         template.HTML
	WordCount       uint64
	TableOfContents Section
	LinesOfCode     uint64
	Tags            []string
	Metadata        Attrs
	Now             time.Time
	Links           Links

	// Parsing
	ParseStart time.Time
	ParseEnd   time.Time

	Filename string
	Raw      []byte // The entire raw file
}

// String returns the Article Title, or the Filename if there is no Title
func (article Article) String() string {
	if article.Title == "" {
		return article.Filename
	}
	return article.Title
}

// SaveAs returns the filename for the output HTML file
func (article Article) SaveAs() string {
	return article.Slug + HTMLExt
}

// RenderWith renders the Article with the given Template
func (article Article) RenderWith(tmpl *template.Template) ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	if err := tmpl.ExecuteTemplate(buffer, "_layout.html", article); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// URL returns the URL of the article
func (article Article) URL() string {
	return article.Links.ForArticle(article)
}

// Transform modifies the given Article, for instance by converting
// its markdown to HTML, or performing syntax highlighting
func (article *Article) Transform(steps ...Transformer) error {
	article.ParseStart = time.Now()
	for _, step := range steps {
		if err := step(article); err != nil {
			return fmt.Errorf(
				`Error while transforming article '%s' with %s: %s`,
				article, step, err,
			)
		}
	}
	article.ParseEnd = time.Now()
	return nil
}

func (article Article) ParseDuration() time.Duration {
	return article.ParseEnd.Sub(article.ParseStart)
}

// NewArticle creates a new article
func NewArticle(filename string) Article {
	article := newArticle()
	article.Filename = filename
	return article
}

func newArticle() Article {
	return Article{
		Metadata: Attrs{},
	}
}
