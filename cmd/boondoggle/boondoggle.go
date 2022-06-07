package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aodin/boondoggle"
	"github.com/fsnotify/fsnotify"
)

var inputDir string
var outputDir string
var templateDir string
var previewCount int
var watch bool

func init() {
	flag.StringVar(&inputDir, "in", ".", "input directory")
	flag.StringVar(&outputDir, "out", "./dist", "input directory")
	flag.StringVar(&templateDir, "tmpl", "", "template directory")
	flag.IntVar(&previewCount, "previews", 4, "number of previews")
	flag.BoolVar(&watch, "watch", false, "watch the input and template directories")
}

func parse() {
	start := time.Now() // Record total time to parse

	// Parse the input directory
	// TODO logging verbosity
	// fmt.Printf("Parsing articles directory '%s'\n", inputDir)

	// TODO need flags for input directory, output directory
	bd, err := boondoggle.ParseDirectory(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	// TODO logging verbosity
	// fmt.Printf("Found %d articles\n", len(bd.Articles))

	// If a template directory was provided, parse templates
	tmpls := boondoggle.Templates{}
	if templateDir != "" {
		// TODO logging verbosity
		// fmt.Printf("Parsing template directory '%s'\n", templateDir)
		if tmpls, err = boondoggle.ParseTemplates(templateDir); err != nil {
			log.Fatalf("Error while parsing templates: %s", err)
		}
		// fmt.Printf("Found %d templates\n", len(tmpls))
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

	// Get the tags template
	tagsTmpl := boondoggle.ExampleTagsTemplate
	if tmpl := tmpls["tags"]; tmpl != nil {
		tagsTmpl = tmpl
	}

	// File flags
	flags := os.O_RDWR + os.O_CREATE + os.O_TRUNC

	// Render the index
	{
		// Preview a few articles
		n := len(bd.Articles)
		if n > previewCount {
			n = previewCount
		}
		previews := bd.Articles[:n]
		attrs := map[string]interface{}{
			"Articles": previews,
			"Now":      bd.BuildTime,
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
			"Now":      bd.BuildTime,
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

		outputPath := filepath.Join(articleDir, article.SaveAs())
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
		bytesWritten := boondoggle.HumanizeBytes(n)
		fmt.Sprintf("%s", bytesWritten)

		// TODO Logging verbosity? Log each article
		// fmt.Printf(
		// 	"Wrote '%s': %s in %s\n",
		// 	article.Title,
		// 	boondoggle.HumanizeBytes(n),
		// 	article.ParseDuration(),
		// )
	}

	tagsDir := filepath.Join(outputDir, "tags")
	if err = os.MkdirAll(tagsDir, 0755); err != nil {
		log.Fatalf("Error while creating tags directory: %s", err)
	}

	// Render the tags index
	if tagsTmpl != nil {
		attrs := map[string]interface{}{
			"Tags": bd.Tags(),
			"Now":  bd.BuildTime,
		}
		indexPath := filepath.Join(tagsDir, "index.html")
		f, err := os.OpenFile(indexPath, flags, 0644)
		if err != nil {
			log.Fatalf("Error while opening file for tags index: %s", err)
		}
		defer f.Close()
		if err := tagsTmpl.Execute(f, attrs); err != nil {
			log.Fatalf("Error while writing tags index: %s", err)
		}
	}

	// Write the tags
	tagTmpl := tmpls["tag"]
	if tagTmpl != nil {
		// TODO logging verbosity
		// fmt.Printf("Writing %d tags\n", len(bd.ByTag))
		for tag, articles := range bd.ByTag {
			attrs := map[string]interface{}{
				"Tag":      tag,
				"Articles": articles,
				"Now":      bd.BuildTime,
			}

			if len(articles) == 1 {
				attrs["Label"] = "Article"
			} else {
				attrs["Label"] = "Articles"
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

	log.Printf(
		"Wrote %d articles in %d ms\n",
		len(bd.Articles),
		time.Since(start).Milliseconds(),
	)
}

func main() {
	flag.Parse()

	if watch {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		parse() // Perform an initial parse

		done := make(chan bool)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					// TODO What events to listen for? CREATE, WRITE
					log.Println("event:", event)
					parse() // Perform a parse whenever a file changes

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		// Watch the input and template directories
		err = watcher.Add(inputDir)
		if err != nil {
			log.Fatal(err)
		}

		err = watcher.Add(templateDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Watching %s and %s for changes\n", inputDir, templateDir)

		<-done

	} else {
		// Run once and exit
		parse()
	}
}
