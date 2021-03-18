package boondoggle

import "fmt"

type Links interface {
	ForArticle(Article) string
}

type UseDefault struct{}

var _ Links = UseDefault{}

func (links UseDefault) ForArticle(article Article) string {
	// TODO Access to article filename?
	return fmt.Sprintf("/articles/%s.html", article.Slug)
}

var _ Links = UseSlugs{}

type UseSlugs struct{}

func (links UseSlugs) ForArticle(article Article) string {
	return fmt.Sprintf("/articles/%s", article.Slug)
}
