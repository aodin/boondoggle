package main

import (
	"log"

	"github.com/aodin/boondoggle/syntax"
)

func main() {
	example := []byte(`[int(x) for x in y]\n`)

	pygmentize := syntax.Pygmentize{}
	out, err := pygmentize.Highlight(example, "py")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", out)
}
