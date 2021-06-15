package boondoggle

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

// ExtractPreview will parse and remove a preview from raw markdown. The
// preview markdown will be converted to HTML after parsing.
func ExtractPreview(article *Article) (err error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	var opened bool
	var preview strings.Builder
	for scanner.Scan() {
		text := scanner.Text()
		// Since the preview may stretch multiple lines, do not stop
		// parsing until the closing comment is found
		if strings.HasPrefix(text, PreviewPrefix) {
			text = strings.TrimPrefix(text, PreviewPrefix)
			opened = true
		}

		if opened {
			if strings.HasSuffix(text, ClosingComment) {
				text = strings.TrimRight(text, ClosingComment)
				opened = false
			} else {
				text += NewLine
			}
			preview.WriteString(text)
		} else {
			if _, err = out.WriteString(text + NewLine); err != nil {
				return
			}
		}
	}

	// If a hard-coded preview was included, parse the markdown to HTML
	if preview.Len() > 0 {
		article.Preview = template.HTML(blackfriday.MarkdownCommon([]byte(preview.String())))
		// Update the raw article to remove the parsed preview
		article.Raw = out.Bytes()
	}
	return
}

// ExtractPreview must have the Transformer function signature
var _ = Transformer(ExtractPreview)

type previewer struct {
	minLength int
	truncate  bool // Do not include opening and closing tag
}

func (pre previewer) Parse(article *Article) error {
	// If the article already has a preview, skip
	if article.Preview != "" {
		return nil
	}

	var preview bytes.Buffer
	tokenizer := html.NewTokenizer(strings.NewReader(string(article.HTML)))
	depth := 0
	previewLength := 0
	alreadyOpenedTag := false

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
					if !(pre.truncate && !alreadyOpenedTag) {
						if _, err := preview.Write(tokenizer.Raw()); err != nil {
							return err
						}
					}
					alreadyOpenedTag = true
				} else {
					depth -= 1
				}

				if tokenType == html.EndTagToken {
					if previewLength > pre.minLength {
						if !pre.truncate {
							if _, err := preview.Write(tokenizer.Raw()); err != nil {
								return err
							}
						}
						// That's enough
						break Previewing
					} else {
						if !(pre.truncate && !alreadyOpenedTag) {
							if _, err := preview.Write(tokenizer.Raw()); err != nil {
								return err
							}
						}
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

func TruncatedTagPreview(minLength int) Transformer {
	pre := previewer{minLength: minLength, truncate: true}
	return pre.Parse
}
