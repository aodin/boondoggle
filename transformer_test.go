package boondoggle

import (
	"strings"
	"testing"
)

func TestTransformer(t *testing.T) {
	name := Transformer(Noop).String()
	if !strings.HasSuffix(name, "Noop") {
		t.Errorf("Unexpected Transformer name: %s", name)
	}
}
