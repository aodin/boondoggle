package boondoggle

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/aodin/boondoggle/syntax"
)

var (
	pre = formatters.Register("pre", html.New(html.WithClasses(true)))
)

func ChromaCode(article *Article) error {
	// Read line by line - saving code blocks
	scanner := bufio.NewScanner(bytes.NewBuffer(article.Raw))

	var b []byte
	out := bytes.NewBuffer(b)

	var code syntax.CodeBlock

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, CodeFence) {
			if code.Exists() {
				// Already in a code block - time to close and parse!
				if err := chromaCodeToHTML(out, code); err != nil {
					return fmt.Errorf("ChromaCode error: %s", err)
				}
				out.WriteString("\n")     // Add a newline
				code = syntax.CodeBlock{} // Clear the existing code block
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
			article.LinesOfCode += 1
			if _, err := code.Block.WriteString(line + NewLine); err != nil {
				return fmt.Errorf("ChromaCode error: %s", err)
			}
			continue
		}

		if _, err := out.WriteString(line + NewLine); err != nil {
			return fmt.Errorf("ChromaCode error: %s", err)
		}
	}

	// If parsing ended before the code block closed, return an error
	if code.Exists() {
		return fmt.Errorf("ChromaCode: EOF while within a code block")
	}

	article.Raw = out.Bytes()
	return nil
}

func chromaCodeToHTML(w io.Writer, code syntax.CodeBlock) error {
	return quick.Highlight(w, code.Block.String(), code.Lang, "pre", "monokai")
}

// ChromaCode must have the Transformer function signature
var _ = Transformer(ChromaCode)
