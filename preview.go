package boondoggle

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type previewer struct {
	minLength int
}

func (pre previewer) Parse(article *Article) error {
	var preview bytes.Buffer
	tokenizer := html.NewTokenizer(strings.NewReader(string(article.HTML)))
	depth := 0
	previewLength := 0

Previewing:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return tokenizer.Err()
		case html.TextToken:
			if depth > 0 {
				n, err := preview.Write(tokenizer.Raw())
				if err == io.EOF {
					break Previewing
				} else if err != nil {
					return nil
				}
				// TODO How to stop mid-paragraph?
				previewLength += n
			}
		case html.StartTagToken, html.EndTagToken:
			name, _ := tokenizer.TagName()
			if len(name) == 1 && name[0] == 'p' {
				if tokenType == html.StartTagToken {
					depth += 1
				} else {
					depth -= 1
				}
				if _, err := preview.Write(tokenizer.Raw()); err != nil {
					return err
				}
				if tokenType == html.EndTagToken {
					if previewLength > pre.minLength {
						// That's enough
						break Previewing
					}
				}

			} else if depth > 0 {
				// TODO Whitelist tags that can be included in previews?
				if _, err := preview.Write(tokenizer.Raw()); err != nil {
					return err
				}
			}
		}
	}
	article.Preview = template.HTML(preview.String())
	return nil
}

var _ = Transformer(previewer{}.Parse)

func Preview(minLength int) Transformer {
	pre := previewer{minLength: minLength}
	return pre.Parse
}
