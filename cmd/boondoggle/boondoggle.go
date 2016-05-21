package main

import (
	"log"

	"github.com/aodin/boondoggle"
)

func main() {
	// TODO need flags for input directory, output directory
	bd, err := boondoggle.ParseDirectory("../../testdata/")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Articles:", len(bd.Articles))
	log.Println("Article", bd.Articles[0].Title)

	out, err := bd.Articles[0].Render(boondoggle.ExampleArticleTemplate)
	if err != nil {
		log.Fatalf("Error while rendering template: %s", err)
	}
	log.Printf("%s\n", out)
}
