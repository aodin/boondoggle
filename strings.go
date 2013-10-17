package boondoggle

import (
	"regexp"
	"strings"
)

var slugifyClean = regexp.MustCompile(`[^\w\s-]`)
var slugifySpace = regexp.MustCompile(`[-\s]+`)

func Slugify(input string) string {
	// Remove anything that isn't a digit, word character or dash
	lowered := string(slugifyClean.ReplaceAll([]byte(input), []byte("")))
	lowered = strings.TrimSpace(strings.ToLower(lowered))
	return string(slugifySpace.ReplaceAll([]byte(lowered), []byte("-")))
}

func UnSnakeCase(input string) string {
	words := strings.Split(input, "_")
	titles := make([]string, len(words))
	// TODO Inefficient method for capitalization, look at unicode.ToTitle
	for index, word := range words {
		titles[index] = strings.Title(word)
	}
	return strings.Join(titles, " ")
}
