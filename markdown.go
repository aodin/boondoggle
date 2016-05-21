package boondoggle

import (
	"html/template"

	"github.com/russross/blackfriday"
)

// MarkdownToHTML will convert the raw markdown bytes to an HTML template.
func MarkdownToHTML(article *Article) (err error) {
	article.HTML = template.HTML(blackfriday.MarkdownCommon(article.Raw))
	return nil
}

// MarkdownToHTML must have the Transformer function signature
var _ = Transformer(MarkdownToHTML)
