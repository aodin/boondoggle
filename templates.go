package boondoggle

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	HTMLExt = ".html"
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

type Templates map[string]*template.Template

func isHTML(name string) bool {
	return strings.ToLower(filepath.Ext(name)) == HTMLExt
}

func isAux(name string) bool {
	return strings.HasPrefix(name, "_")
}

func ParseTemplates(path string) (Templates, error) {
	parsed := Templates{}

	// TODO walk?
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return parsed, err
	}

	// Any files prefixed with '_' will be added to each template
	aux := []string{}
	for _, file := range files {
		if isHTML(file.Name()) && isAux(file.Name()) {
			aux = append(aux, filepath.Join(path, file.Name()))
		}
	}

	for _, file := range files {
		if !isHTML(file.Name()) || isAux(file.Name()) {
			continue
		}

		fullpath := filepath.Join(path, file.Name())
		name := strings.TrimSuffix(strings.ToLower(file.Name()), HTMLExt)

		if parsed[name], err = template.ParseFiles(append(aux, fullpath)...); err != nil {
			return parsed, err
		}
	}
	return parsed, nil
}
