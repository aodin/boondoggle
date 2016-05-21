package syntax

// noop implements the Highlighter interface but does nothing
type noop struct{}

var _ Highlighter = noop{}

func (h noop) Highlight(text []byte, lang string) ([]byte, error) {
	return text, nil
}
