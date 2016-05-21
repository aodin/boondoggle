package boondoggle

import (
	"strings"
	"time"
)

const ISO8601Date = "2006-01-02"

// ParseFilename will parse the filename for the slug and date.
// The filename must be in the format YYYY-MM-DD_title.md
func ParseFilename(article *Article) (err error) {
	// TODO case insensitive trimming?
	filename := strings.TrimSuffix(article.Filename, ".md")
	filename = strings.TrimSuffix(filename, ".MD")

	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		date, err := time.Parse(ISO8601Date, parts[0])
		if err == nil {
			article.Date = date
			article.Title = strings.Join(parts[1:], "_")
			article.Slug = Slug(article.Title)
		} else {
			article.Title = strings.Join(parts, "_")
			article.Slug = Slug(article.Title)
		}
	} else {
		article.Filename = parts[0]
		article.Slug = Slug(article.Title)
	}

	return nil
}

// ParseFilename must have the Transformer function signature
var _ = Transformer(ParseFilename)
