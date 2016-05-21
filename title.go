package boondoggle

import (
	"bufio"
	"bytes"
	"strings"
)

// TODO Is there a better way to get the unread scanner bytes?
func readRemainder(scanner *bufio.Scanner) ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	for scanner.Scan() {
		if _, err := buffer.Write(scanner.Bytes()); err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(NewLine); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// ExtractTitle will parse and remove an atx or setext H1 title from the first
// line (and second if setext) of the markdown file.
// TODO ExtractTitle will add a newline - it shouldn'y
func ExtractTitle(article *Article) (err error) {
	buffer := bytes.NewBuffer(article.Raw)
	scanner := bufio.NewScanner(buffer)

	// Read the first line
	if !scanner.Scan() {
		// No text!
		return
	}

	first := scanner.Text()
	if strings.HasPrefix(first, TitleAtx) {
		if len(first) > 1 && string(first[1]) != TitleAtx {
			article.Title = strings.TrimSpace(first[1:])

			// Stop here - Ignore a setext if we already have atx
			// Dump the remainder of the buffer into raw
			// article.Raw = buffer.Bytes() // TODO Doesn't work
			article.Raw, err = readRemainder(scanner)
			return
		} else {
			// Don't even bother with setext - the atx is funky
			return
		}
	}

	// Read the second line
	if !scanner.Scan() {
		// No more text - don't bother with setext
		return
	}

	second := scanner.Text()
	if len(strings.TrimRight(second, Space)) == strings.Count(second, TitleSetext) {
		// Settext!
		article.Title = strings.TrimSpace(first)

		// Dump the remainder of the buffer into raw
		// article.Raw = buffer.Bytes() // TODO Doesn't work
		article.Raw, err = readRemainder(scanner)
		return
	}

	// No title was found - do nothing to the raw article
	return
}

// ExtractTitle must have the Transformer function signature
var _ = Transformer(ExtractTitle)
