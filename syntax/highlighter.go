package syntax

// Highlighter is the interface syntax highlighters must implement
type Highlighter interface {
	Highlight(text []byte, lang string) ([]byte, error)
}
