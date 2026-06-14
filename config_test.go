package boondoggle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Input != "." {
		t.Errorf("Unexpected default input: %q", config.Input)
	}
	if config.Output != "./dist" {
		t.Errorf("Unexpected default output: %q", config.Output)
	}
	if config.Previews != 4 {
		t.Errorf("Unexpected default previews: %d", config.Previews)
	}
}

func TestLoadConfig(t *testing.T) {
	contents := `
input = "articles"
output = "public"
templates = "templates"
previews = 10

[feed]
title = "Example Blog"
url = "https://example.com"
description = "An example blog"
author = "Jane Doe"
email = "jane@example.com"
items = 20
`
	path := filepath.Join(t.TempDir(), "boondoggle.toml")
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("Failed to write config: %s", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %s", err)
	}

	if config.Input != "articles" {
		t.Errorf("Unexpected input: %q", config.Input)
	}
	if config.Output != "public" {
		t.Errorf("Unexpected output: %q", config.Output)
	}
	if config.Templates != "templates" {
		t.Errorf("Unexpected templates: %q", config.Templates)
	}
	if config.Previews != 10 {
		t.Errorf("Unexpected previews: %d", config.Previews)
	}

	feed := config.Feed
	if feed.Title != "Example Blog" {
		t.Errorf("Unexpected feed title: %q", feed.Title)
	}
	if feed.Link != "https://example.com" {
		t.Errorf("Unexpected feed url: %q", feed.Link)
	}
	if feed.Description != "An example blog" {
		t.Errorf("Unexpected feed description: %q", feed.Description)
	}
	if feed.Author != "Jane Doe" {
		t.Errorf("Unexpected feed author: %q", feed.Author)
	}
	if feed.Email != "jane@example.com" {
		t.Errorf("Unexpected feed email: %q", feed.Email)
	}
	if feed.Limit != 20 {
		t.Errorf("Unexpected feed items: %d", feed.Limit)
	}
}

// Values omitted from the file should keep their defaults.
func TestLoadConfigPartial(t *testing.T) {
	path := filepath.Join(t.TempDir(), "boondoggle.toml")
	if err := os.WriteFile(path, []byte("output = \"build\"\n"), 0644); err != nil {
		t.Fatalf("Failed to write config: %s", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %s", err)
	}
	if config.Output != "build" {
		t.Errorf("Unexpected output: %q", config.Output)
	}
	// Untouched fields retain their defaults.
	if config.Input != "." {
		t.Errorf("Expected default input, got %q", config.Input)
	}
	if config.Previews != 4 {
		t.Errorf("Expected default previews, got %d", config.Previews)
	}
}
