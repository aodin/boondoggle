package boondoggle

import (
	"bufio"
	"bytes"
	"strings"
)

// ExtractTags will parse and remove the tags from the raw markdown. The
// Tags will be converted to slugs after parsing.
func ExtractTags(article *Article) (err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, TagsPrefix) {
			text = strings.TrimPrefix(text, TagsPrefix)
			text = strings.TrimRight(text, ClosingComment)
			article.Tags = normalizeTags(strings.Split(text, ","))
		} else {
			if _, err = out.WriteString(text + NewLine); err != nil {
				return
			}
		}
	}

	article.Raw = out.Bytes()
	return
}

func normalizeTags(in []string) (out []string) {
	for _, tag := range in {
		if tag = Slug(tag); tag != "" {
			out = append(out, tag)
		}
	}
	return
}

// ExtractTags must have the Transformer function signature
var _ = Transformer(ExtractTags)
