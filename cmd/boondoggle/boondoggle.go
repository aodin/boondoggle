package main

import (
	"log"

	"github.com/aodin/boondoggle"
)

func main() {
	// Example markdown
	content := []byte(`Example Article
=======

<!-- tags: SQL, python -->

I am a paragraph. I contain words.`)

	// ugh
	content = append(content, []byte("```sql\n")...)
	content = append(content, []byte(`
package main

// Comment!

func main() {
	println("Hello, world!")
}
`)...)
	content = append(content, []byte("```\n")...)

	// Call boondoggle with the following Transformers
	pipeline := []boondoggle.Transformer{
		boondoggle.ExtractTags,
		boondoggle.ExtractTitle,
	}

	// TODO need flags for input directory, output directory
	article, err := boondoggle.ParseMarkdown(content, pipeline...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(article.Title)
	log.Println(article.Tags)
	log.Println(article.HTML)
}
