package boondoggle

import (
	"regexp"
	"strings"
	"unicode"
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
	// TODO There is also a unicode method for IsTitle
	capitalize := true
	if len(words) > 0 {
		firstWord := []rune(words[0])
		if len(firstWord) > 0 && unicode.IsUpper(firstWord[0]) {
			capitalize = false
		}
	}
	for index, word := range words {
		if capitalize {
			titles[index] = strings.Title(word)
		} else {
			titles[index] = word
		}
	}
	return strings.Join(titles, " ")
}

var filenameSplitter = regexp.MustCompile(`(?P<date>\d+-\d+-\d+)[ -_](?P<title>.*)`)

// Split a date and title from the given input
func SplitFilename(input string) (string, string) {
	// TODO Or could use the split functionality of regexp
	results := filenameSplitter.FindStringSubmatch(input)
	if results == nil || len(results) < 3 {
		return "", ""
	}
	return results[1], results[2]
}
