package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/russross/blackfriday"

	"github.com/aodin/boondoggle/syntax"
)

type CodeBlock struct {
	Lang  string
	Block *bytes.Buffer
}

func (code CodeBlock) Exists() bool {
	return code.Block != nil
}

func (code CodeBlock) Content() ([]byte, string) {
	return code.Block.Bytes(), code.Lang
}

func NewCodeBlock() CodeBlock {
	var b []byte
	return CodeBlock{Block: bytes.NewBuffer(b)}
}

const codeFence = "```"

// TODO parser with composable Highlighter

func MarkdownWithCode(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read line by line - saving code blocks
	scanner := bufio.NewScanner(file) // Defaults to ScanLines

	var b []byte
	parsed := bytes.NewBuffer(b)

	// TODO Is there a smarter way to parse this?
	var code CodeBlock

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, codeFence) {
			if code.Exists() {
				// Already in a code block - time to close and parse!
				htmlCode, err := ReplaceCode(code.Content())
				if err != nil {
					return nil, fmt.Errorf("ReplaceCode error: %s", err)
				}
				// TODO ignore error
				parsed.Write(htmlCode)

				code = CodeBlock{} // TODO Method to end the code block?
				continue
			} else {
				// Create a new code block and continue
				code = NewCodeBlock()
				code.Lang = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, codeFence)))
				continue
			}
		}

		if code.Exists() {
			code.Block.WriteString(line)
			code.Block.WriteRune('\n')
			continue
		}

		// TODO there are ignored errors
		parsed.WriteString(line)
		parsed.WriteRune('\n')
	}

	// If parsing ended before the code block closed, return an error
	if code.Exists() {
		return nil, fmt.Errorf("encountered EOF while within a code block")
	}

	return blackfriday.MarkdownCommon(parsed.Bytes()), nil
}

// ReplaceCode replaces all ``` fence blocks with highlighted code
// If there is a string after the opening ```, it will be passed
// as the language option to the highlighter
func ReplaceCode(text []byte, lang string) ([]byte, error) {
	pygmentize := syntax.Pygmentize{}
	return pygmentize.Highlight(text, lang)
}

func main() {
	// Import a file - replace its code blocks with highlighted HTML
	out, err := MarkdownWithCode("../../testdata/third.md")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", out)

	// Inline example
	// example := []byte(`[int(x) for x in y]\n`)
	// pygmentize := syntax.Pygmentize{}
	// out, err := pygmentize.Highlight(example, "py")
	// if err != nil {
	//  log.Fatal(err)
	// }
	// log.Printf("%s", out)
}
