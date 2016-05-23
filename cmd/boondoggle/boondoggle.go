package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aodin/boondoggle"
)

var inputDir string
var outputDir string
var templateDir string

func init() {
	flag.StringVar(&inputDir, "in", ".", "input directory")
	flag.StringVar(&outputDir, "out", "./dist", "input directory")
	flag.StringVar(&templateDir, "tmpl", "", "template directory")
}

func main() {
	flag.Parse()

	// Parse the input directory
	log.Printf("Parsing input directory '%s'...", inputDir)

	// TODO need flags for input directory, output directory
	bd, err := boondoggle.ParseDirectory(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Parsed %d articles", len(bd.Articles))

	// If a template directory was provided, parse templates
	tmpls := boondoggle.Templates{}
	if templateDir != "" {
		log.Printf("Parsing template directory '%s'...", templateDir)
		if tmpls, err = boondoggle.ParseTemplates(templateDir); err != nil {
			log.Fatalf("Error while parsing templates: %s", err)
		}
		log.Printf("Parsed %d templates", len(tmpls))
	}

	// Does the destination directory exist?
	articleDir := filepath.Join(outputDir, "articles")
	if err = os.MkdirAll(articleDir, 0755); err != nil {
		log.Fatalf("Error while creating output directory: %s", err)
	}

	// Get the article template, if one was not parsed, use the example
	articleTmpl := boondoggle.ExampleArticleTemplate
	if tmpl := tmpls["article"]; tmpl != nil {
		articleTmpl = tmpl
	}

	// Get the index template
	indexTmpl := boondoggle.ExampleIndexTemplate
	if tmpl := tmpls["index"]; tmpl != nil {
		indexTmpl = tmpl
	}

	// File flags
	flags := os.O_RDWR + os.O_CREATE + os.O_TRUNC

	// Render the index
	{
		// Preview a few articles
		n := len(bd.Articles)
		if n > 4 {
			n = 4
		}
		previews := bd.Articles[:n]
		attrs := map[string]interface{}{
			"Articles": previews,
		}

		indexPath := filepath.Join(outputDir, "index.html")
		f, err := os.OpenFile(indexPath, flags, 0644)
		if err != nil {
			log.Fatalf("Error while opening file for index: %s", err)
		}
		defer f.Close()
		if err := indexTmpl.Execute(f, attrs); err != nil {
			log.Fatalf("Error while writing index: %s", err)
		}
	}

	// Render the Articles index
	articlesTmpl := tmpls["articles"]
	if articlesTmpl != nil {
		attrs := map[string]interface{}{
			"Articles": bd.Articles,
		}

		indexPath := filepath.Join(articleDir, "index.html")
		f, err := os.OpenFile(indexPath, flags, 0644)
		if err != nil {
			log.Fatalf("Error while opening file for index: %s", err)
		}
		defer f.Close()
		if err := articlesTmpl.Execute(f, attrs); err != nil {
			log.Fatalf("Error while writing articles index: %s", err)
		}
	}

	// Render each article
	for _, article := range bd.Articles {
		// TODO Render directory to the given file handler
		out, err := article.RenderWith(articleTmpl)
		if err != nil {
			log.Fatalf(
				"Error while rendering template for %s: %s",
				article.Slug, err,
			)
		}

		outputPath := filepath.Join(articleDir, article.Slug+".html")
		f, err := os.OpenFile(outputPath, flags, 0644)
		if err != nil {
			log.Fatalf(
				"Error while opening file for %s: %s",
				article.Slug, err,
			)
		}
		defer f.Close()

		n, err := f.Write(out)
		if err != nil {
			log.Fatalf(
				"Error while writing file for %s: %s",
				article.Slug, err,
			)
		}
		log.Printf("Wrote %d bytes for %s", n, article.Slug)
	}

	tagsDir := filepath.Join(outputDir, "tags")
	if err = os.MkdirAll(tagsDir, 0755); err != nil {
		log.Fatalf("Error while creating tags directory: %s", err)
	}

	// Write the tags
	tagTmpl := tmpls["tag"]
	if tagTmpl != nil {
		log.Printf("Writing %d tags...", len(bd.ByTag))
		for tag, articles := range bd.ByTag {
			attrs := map[string]interface{}{
				"Tag":      tag,
				"Articles": articles,
			}
			outputPath := filepath.Join(tagsDir, tag+".html")
			f, err := os.OpenFile(outputPath, flags, 0644)
			if err != nil {
				log.Fatalf("Error while opening file for tag %s: %s", tag, err)
			}
			defer f.Close()
			if err := tagTmpl.Execute(f, attrs); err != nil {
				log.Fatalf("Error while writing tag %s: %s", tag, err)
			}
		}
	}
}
