// Boondoggle is a static site generator written in Go.

package boondoggle

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	MarkdownExt = ".md" // MarkdownExt is the common markdown file ending
	HTMLExt     = ".html"
)

var DefaultProcessor = []Transformer{
	ParseFilename,
	ExtractFrontMatter,
	ExtractTitle,
	ExtractTags,
	ExtractPreview,
	ChromaCode,
	MarkdownToHTML,
	TruncatedTagPreview(200),
}

// Boondoggle builds .html files from a directory of markdown files.
type Boondoggle struct {
	Links    Links // Writes URLs
	Articles Articles

	ByTitle map[string]Article
	ByTag   map[string]Articles

	Metadata  Attrs
	BuildTime time.Time
}

// Tags returns tags in alphabetical order
func (bd Boondoggle) Tags() (tags []string) {
	for tag, _ := range bd.ByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return
}

func (bd *Boondoggle) ReadDirectory(path string) error {
	// Parse each file in the directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		extension := strings.ToLower(filepath.Ext(name))

		if extension != MarkdownExt {
			continue
		}

		fullpath := filepath.Join(path, name)

		// TODO or use io.Reader?
		content, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

		article := NewArticle(name)
		article.Links = bd.Links
		article.Raw = content
		article.Now = bd.BuildTime
		bd.Articles = append(bd.Articles, article)
	}
	return nil
}

// New creates a new Boondoggle. The New method does not need to be used
// directly - use ParseDirectory instead
func New() *Boondoggle {
	return &Boondoggle{
		Links:     UseSlugs{},
		ByTitle:   make(map[string]Article),
		ByTag:     make(map[string]Articles),
		Metadata:  Attrs{},
		BuildTime: time.Now(),
	}
}

// ParseDirectory will parse all markdown files in the given directory
func ParseDirectory(path string, steps ...Transformer) (*Boondoggle, error) {
	bd := New()
	if err := bd.ReadDirectory(path); err != nil {
		return nil, err
	}

	// If no transformers were provided, use the default processor
	if len(steps) == 0 {
		steps = DefaultProcessor
	}

	for index, article := range bd.Articles {
		if err := article.Transform(steps...); err != nil {
			return nil, err
		}
		bd.Articles[index] = article

		// Aggregate tags
		for _, tag := range article.Tags {
			bd.ByTag[tag] = append(bd.ByTag[tag], article)
		}
	}

	// Sort articles
	bd.Articles.SortMostRecentArticlesFirst()
	for tag := range bd.ByTag {
		bd.ByTag[tag].SortMostRecentArticlesFirst()
	}
	return bd, nil
}
