/*
Boondoggle is a static site generator written in Go.
*/
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
)

// Boondoggle builds .HTML files from a directory of markdown files.
type Boondoggle struct {
	Articles Articles

	ByTitle map[string]Article  // Includes path - TODO pointer?
	ByTag   map[string]Articles // TODO pointer?

	// TODO Metadata?
	BuildTime time.Time
}

// TODO Tags returns tags in alphabetical order
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
		article.Raw = content
		bd.Articles = append(bd.Articles, article)
	}
	return nil
}

// New creates a new Boondoggle. The New method does not need to be used
// directly - use ParseDirectory instead
func New() *Boondoggle {
	return &Boondoggle{
		ByTitle:   make(map[string]Article),
		ByTag:     make(map[string]Articles),
		BuildTime: time.Now(),
	}
}

// ParseDirectory will parse all markdown files in the given directory
// TODO Walk the entire directory structure?
// TODO noop HTML files?
func ParseDirectory(path string, steps ...Transformer) (*Boondoggle, error) {
	bd := New()
	if err := bd.ReadDirectory(path); err != nil {
		return nil, err
	}

	// For each article, perform the default actions, unless alternative
	// transformers have been given
	if len(steps) == 0 {
		steps = []Transformer{
			ExtractTitle,
			ExtractTags,
			PygmentizeCode,
			MarkdownToHTML,
		}
	}

	for _, article := range bd.Articles {
		for _, step := range steps {
			if err := step(&article); err != nil {
				return nil, err
			}
		}
	}

	return bd, nil
}
