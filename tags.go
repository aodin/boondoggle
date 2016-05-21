package boondoggle

import (
	"bufio"
	"bytes"
	"strings"
)

// TODO Rather unfriendly
const TagsPrefix = "<!-- tags:"

// ExtractTags will parse and remove the tags from the raw markdown. The
// Tags will be slugified after parsing.
func ExtractTags(article *Article) (err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, TagsPrefix) {
			text = strings.TrimPrefix(text, TagsPrefix)
			text = strings.TrimRight(text, "-> ")

			tags := strings.Split(text, ",")
			for _, tag := range tags {
				// TODO Allow unicode tags?
				if tag = Slug(tag); tag != "" {
					article.Tags = append(article.Tags, tag)
				}
			}
		} else {
			if _, err = out.WriteString(text + NewLine); err != nil {
				return
			}
		}
	}

	article.Raw = out.Bytes()
	return
}

// ExtractTags must have the Transformer function signature
var _ = Transformer(ExtractTags)
