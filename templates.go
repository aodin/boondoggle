package boondoggle

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Example Templates
var ExampleArticleTemplate = template.Must(template.New("article").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
  </head>
  <body>
    <h1>{{ .Title }}</h1>
    <h4>{{ .Date.Format "Monday, January 2, 2006" }}<h4>
    {{ .HTML }}
  </body>
</html>
`))

var ExampleIndexTemplate = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>My Site</title>
  </head>
  <body>
    {{ range $article := .Articles }}
    <article>
      <h1>{{ $article.Title }}</h1>
    </article>
    {{ end }}
  </body>
    }
</html>
`))

var ExampleTagsTemplate = template.Must(template.New("tags").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>Tags</title>
  </head>
  <body>
    <ul>
    {{ range $tag := .Tags }}
      <li>{{ $tag }}</li>
	  {{ end }}
	</ul>
  </body>
</html>
`))

type Templates map[string]*template.Template

func ParseTemplates(path string) (Templates, error) {
	parsed := Templates{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return parsed, err
	}

	layout := filepath.Join(path, "_layout.html")

	for _, file := range files {
		filename := strings.ToLower(file.Name())
		extension := filepath.Ext(filename)

		if extension != HTMLExt {
			continue
		}

		fullpath := filepath.Join(path, filename)
		name := strings.TrimSuffix(filename, HTMLExt)

		if parsed[name], err = template.ParseFiles(layout, fullpath); err != nil {
			return parsed, err
		}
	}
	return parsed, nil
}
