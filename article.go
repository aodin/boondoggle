package boondoggle

import (
	"bytes"
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
	WordCount       uint64
	TableOfContents Section
	LinesOfCode     uint64
	Tags            []string
	Metadata        map[string]interface{}

	// TODO need methods to create buffers/scanners and reset raw
	Filename string
	Raw      []byte // The entire raw file - TODO use io.Reader?
}

// String returns the Article Title, falling back to Filename if Title does
// not exist
func (article Article) String() string {
	if article.Title == "" {
		return article.Filename
	}
	return article.Title
}

// Render renders the Article with the given Template
func (article Article) Render(tmpl *template.Template) ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	if err := tmpl.Execute(buffer, article); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ParseMarkdown creates an Article from the given markdown,
// optionally running it through any given transformers.
func ParseMarkdown(markdown []byte, pipeline ...Transformer) (Article, error) {
	article := Article{
		Raw:      markdown,
		Metadata: make(map[string]interface{}),
	}

	// Always call MarkdownToHTML at the end
	// TODO prevent MarkdownToHTML from being run multiple times?
	for _, step := range append(pipeline, MarkdownToHTML) {
		if err := step(&article); err != nil {
			return article, err
		}
	}
	return article, nil
}

// NewArticle creates a new article
func NewArticle(filename string) Article {
	return Article{
		Metadata: make(map[string]interface{}),
		Filename: filename,
	}
}
