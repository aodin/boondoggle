package boondoggle

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// ExtractFrontMatter will parse and remove any Front Matter
// (https://jekyllrb.com/docs/frontmatter/) metadata tags from the
// raw markdown. For front matter to be parsed, the first line
// of the markdown must start with "---"
func ExtractFrontMatter(article *Article) (err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	if !scanner.Scan() {
		// No content!
		return
	}

	first := scanner.Text()
	if first != FrontMatterBlock {
		// No front matter!
		return
	}

	var metadata []byte
	matter := bytes.NewBuffer(metadata)

	// Parse all content until another FrontMatterBlock is found
	inBlock := true
	for scanner.Scan() {
		text := scanner.Text()
		if inBlock && text == FrontMatterBlock {
			inBlock = false
		} else if inBlock {
			if _, err = matter.WriteString(text + NewLine); err != nil {
				return
			}
		} else {
			if _, err = out.WriteString(text + NewLine); err != nil {
				return
			}
		}
	}

	if inBlock {
		return fmt.Errorf("Front Matter block was never closed")
	}

	// Parse the entire front matter block
	if err := yaml.Unmarshal(matter.Bytes(), &article.Metadata); err != nil {
		return err
	}

	if title := getTitle(article.Metadata); title != "" {
		article.Title = title
	}
	if tags := getTags(article.Metadata); len(tags) > 0 {
		article.Tags = tags
	}

	article.Raw = out.Bytes()
	return
}

func getTitle(attrs Attrs) string {
	title, _ := attrs["title"].(string)
	return title
}

func getTags(attrs Attrs) []string {
	// tags can either be a string or a list of strings
	tags, ok := attrs["tags"].([]string)
	if !ok {
		if raw, ok := attrs["tags"].(string); ok {
			tags = strings.Split(raw, ", ")
		}
	}
	return normalizeTags(tags)
}

// ExtractFrontMatter must have the Transformer function signature
var _ = Transformer(ExtractFrontMatter)
