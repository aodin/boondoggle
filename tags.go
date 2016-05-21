package boondoggle

import (
	"bufio"
	"bytes"
	"strings"
)

// TODO Rather unfriendly
const TagsPrefix = "[//]: # ("

// ParseTags will parse the tags from the raw markdown. The tags line
// must be present within the first 5 lines of the markdown file.
// Tags will be slugified after parsing. Only one line of tags is allowed.
// The tags line will not be removed since it is a markdown "comment":
// http://stackoverflow.com/a/20885980
// TODO For cross platform support, it should probably be removed
func ParseTags(article *Article) error {
	buffer := bytes.NewBuffer(article.Raw)
	scanner := bufio.NewScanner(buffer)

	n := 0
	for scanner.Scan() && n < 5 {
		text := scanner.Text()
		if strings.HasPrefix(text, TagsPrefix) {
			text = strings.TrimPrefix(text, TagsPrefix)
			text = strings.TrimRight(text, ") ")

			tags := strings.Split(text, ",")
			for _, tag := range tags {
				// TODO Allow unicode tags?
				if tag = Slug(tag); tag != "" {
					article.Tags = append(article.Tags, tag)
				}
			}
			break
		}
		n += 1
	}
	return nil
}

// ParseTags must have the Transformer function signature
var _ = Transformer(ParseTags)
