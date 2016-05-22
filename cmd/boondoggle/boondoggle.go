package main

import (
	"flag"
	"log"

	"github.com/aodin/boondoggle"
)

var input string
var output string

func init() {
	flag.StringVar(&input, "in", ".", "input directory")
	flag.StringVar(&output, "out", "./dist", "input directory")
}

func main() {
	flag.Parse()

	// Parse the input directory
	log.Printf("Parsing input directory '%s'...", input)

	// TODO need flags for input directory, output directory
	bd, err := boondoggle.ParseDirectory(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed %d articles", len(bd.Articles))
	for _, article := range bd.Articles {
		out, err := article.RenderWith(boondoggle.ExampleArticleTemplate)
		if err != nil {
			log.Fatalf("Error while rendering template: %s", err)
		}
		log.Printf("%s\n", out)
	}
}
