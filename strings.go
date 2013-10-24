package boondoggle

import (
	"fmt"
	"regexp"
	"strings"
	"time"
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

func ParseDate(input string) (time.Time, error) {
	return time.Parse("2006-01-02", input)
}

func OutputDate(input time.Time) string {
	return fmt.Sprintf("%s, %s %s", input.Weekday(), input.Month(), input.Format("2, 2006"))
}
