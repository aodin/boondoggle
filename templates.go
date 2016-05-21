package boondoggle

import "html/template"

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
