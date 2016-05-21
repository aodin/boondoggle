package boondoggle

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/aodin/boondoggle/syntax"
)

func PygmentizeCode(article *Article) error {
	// Read line by line - saving code blocks
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	// TODO Is there a smarter way to parse this?
	var code syntax.CodeBlock

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, CodeFence) {
			if code.Exists() {
				// Already in a code block - time to close and parse!
				htmlCode, err := pygmentizeCode(code.Content())
				if err != nil {
					return fmt.Errorf("PygmentizeCode error: %s", err)
				}
				if _, err = out.Write(htmlCode); err != nil {
					return fmt.Errorf("PygmentizeCode error: %s", err)
				}

				code = syntax.CodeBlock{} // TODO Method to end the code block?
				continue
			} else {
				// Create a new code block and continue
				code = syntax.NewCodeBlock()
				code.Lang = strings.ToLower(
					strings.TrimSpace(strings.TrimPrefix(line, CodeFence)),
				)
				continue
			}
		}

		if code.Exists() {
			if _, err := code.Block.WriteString(line + NewLine); err != nil {
				return fmt.Errorf("PygmentizeCode error: %s", err)
			}
			continue
		}

		if _, err := out.WriteString(line + NewLine); err != nil {
			return fmt.Errorf("PygmentizeCode error: %s", err)
		}
	}

	// If parsing ended before the code block closed, return an error
	if code.Exists() {
		return fmt.Errorf("PygmentizeCode: EOF while within a code block")
	}

	article.Raw = out.Bytes()
	return nil
}

func pygmentizeCode(text []byte, lang string) ([]byte, error) {
	pygmentize := syntax.Pygmentize{}
	return pygmentize.Highlight(text, lang)
}

// PygmentizeCode must have the Transformer function signature
var _ = Transformer(PygmentizeCode)
