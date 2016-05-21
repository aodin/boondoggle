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
}
