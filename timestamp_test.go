package boondoggle

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	input := "2013-10-22"
	parsed, err := CreateTimestamp(input)
	if err != nil {
		t.Error("Error during timestamp creation:", err)
	}
	actual := Timestamp{time.Date(2013, time.October, 22, 0, 0, 0, 0, time.UTC)}
	if parsed != actual {
		t.Errorf("Unexpected date: %s != %s", parsed.String(), actual.String())
	}
	// Test the String() method
	output := `Tuesday, October 22, 2013`
	if actual.String() != output {
		t.Errorf("Unexpected output: %s != %s", actual.String(), output)
	}
}
