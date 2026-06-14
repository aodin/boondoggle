package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aodin/boondoggle"
	"github.com/fsnotify/fsnotify"
)

var configPath string
var inputDir string
var outputDir string
var templateDir string
var previewCount int
var watch bool

// Feed metadata
var siteTitle string
var siteURL string
var siteDescription string
var siteAuthor string
var siteEmail string
var feedItems int

func init() {
	defaults := boondoggle.DefaultConfig()

	flag.StringVar(&configPath, "config", "boondoggle.toml", "path to the configuration file")
	flag.StringVar(&inputDir, "in", defaults.Input, "input directory")
	flag.StringVar(&outputDir, "out", defaults.Output, "output directory")
	flag.StringVar(&templateDir, "tmpl", defaults.Templates, "template directory")
	flag.IntVar(&previewCount, "previews", defaults.Previews, "number of previews")
	flag.BoolVar(&watch, "watch", false, "watch the input and template directories")

	flag.StringVar(&siteTitle, "title", "", "site title for the RSS and Atom feeds")
	flag.StringVar(&siteURL, "url", "", "absolute base URL of the site, e.g. https://example.com (enables feeds)")
	flag.StringVar(&siteDescription, "desc", "", "site description for the RSS and Atom feeds")
	flag.StringVar(&siteAuthor, "author", "", "default feed author name")
	flag.StringVar(&siteEmail, "email", "", "default feed author email")
	flag.IntVar(&feedItems, "feeditems", 0, "maximum number of articles in the feeds, set 0 for all")
}

// loadSettings merges the configuration file with the command-line flags.
// Explicit flags take precedence over the config file, which in turn takes
// precedence over the built-in defaults.
func loadSettings() {
	// Track which flags were explicitly set on the command line.
	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })

	config := boondoggle.DefaultConfig()
	if _, err := os.Stat(configPath); err == nil {
		if config, err = boondoggle.LoadConfig(configPath); err != nil {
			log.Fatalf("Error while loading config %s: %s", configPath, err)
		}
	} else if set["config"] {
		// An explicitly requested config file that is missing is an error.
		log.Fatalf("Config file %s does not exist", configPath)
	}

	if !set["in"] {
		inputDir = config.Input
	}
	if !set["out"] {
		outputDir = config.Output
	}
	if !set["tmpl"] {
		templateDir = config.Templates
	}
	if !set["previews"] {
		previewCount = config.Previews
	}
	if !set["title"] {
		siteTitle = config.Feed.Title
	}
	if !set["url"] {
		siteURL = config.Feed.Link
	}
	if !set["desc"] {
		siteDescription = config.Feed.Description
	}
	if !set["author"] {
		siteAuthor = config.Feed.Author
	}
	if !set["email"] {
		siteEmail = config.Feed.Email
	}
	if !set["feeditems"] {
		feedItems = config.Feed.Limit
	}
}

// underOutput reports whether path is the output directory or a file within it.
// It is used to ignore the generator's own writes while watching.
func underOutput(path string) bool {
	out, err := filepath.Abs(outputDir)
	if err != nil {
		return false
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(out, abs)
	if err != nil {
		return false
	}
	// rel does not escape the output directory (no leading "..").
	return rel == "." || !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".."
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

	// Write the RSS and Atom feeds. Feeds require absolute URLs, so they are
	// only generated when a base site URL is provided.
	if siteURL != "" {
		feed := boondoggle.Feed{
			Title:       siteTitle,
			Link:        siteURL,
			Description: siteDescription,
			Author:      siteAuthor,
			Email:       siteEmail,
			Limit:       feedItems,
		}

		feeds := map[string]boondoggle.FeedTransformer{
			"feed.xml": feed.RSS,
			"atom.xml": feed.Atom,
		}
		for name, render := range feeds {
			out, err := render(bd)
			if err != nil {
				log.Fatalf("Error while rendering %s: %s", name, err)
			}

			outputPath := filepath.Join(outputDir, name)
			f, err := os.OpenFile(outputPath, flags, 0644)
			if err != nil {
				log.Fatalf("Error while opening file for %s: %s", name, err)
			}
			defer f.Close()
			if _, err := f.Write(out); err != nil {
				log.Fatalf("Error while writing %s: %s", name, err)
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
	loadSettings()

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
					// Ignore the files we generate; otherwise writing the
					// output would trigger an endless rebuild loop.
					if underOutput(event.Name) {
						continue
					}
					// TODO What events to listen for? CREATE, WRITE
					log.Println("event:", event)
					// Reload settings so config file edits take effect, then
					// rebuild.
					loadSettings()
					parse()

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		// Watch the input directory, and the template and config files when
		// they are present.
		watched := []string{inputDir}
		if templateDir != "" {
			watched = append(watched, templateDir)
		}
		if _, err := os.Stat(configPath); err == nil {
			watched = append(watched, configPath)
		}
		for _, path := range watched {
			if err := watcher.Add(path); err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("Watching %s for changes\n", strings.Join(watched, ", "))

		<-done

	} else {
		// Run once and exit
		parse()
	}
}
