package syntax

// Highlighter is the interface that syntax highlighters must implement
type Highlighter interface {
	Highlight(text []byte, lang string) ([]byte, error)
}
