/*
Package syntax contains a standard Highlighter interface and multiple
implementations.
*/
package syntax

import "bytes"

// CodeBlock is a ``` fenced block of code that will be syntax highlighted
type CodeBlock struct {
	Lang  string
	Block *bytes.Buffer
	// TODO count lines of code?
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
