package boondoggle

import (
	"html/template"
	"testing"
)

var mdPreview = `
# Title

### Subtitle

I am just some text. Maybe *with* styles or [link](/somewhere).

That's all. Long content. More rambling. blah blah blah blah blah.

Another paragraph that should be excluded.
`

var expectedPreview = template.HTML(`I am just some text. Maybe <em>with</em> styles or <a href="/somewhere">link</a>.</p><p>That&rsquo;s all. Long content. More rambling. blah blah blah blah blah.`)

func TestPreview(t *testing.T) {
	var mock Article
	mock.Raw = []byte(mdPreview)

	if err := MarkdownToHTML(&mock); err != nil {
		t.Fatal(err)
	}

	if err := TruncatedTagPreview(64)(&mock); err != nil {
		t.Fatal(err)
	}

	if mock.Preview != expectedPreview {
		t.Errorf("unexpected preview: %s", mock.Preview)
	}
}
