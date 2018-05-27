package views

import (
	"html/template"

	"hawx.me/code/papiermache/data"
)

const read = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>{{ .Title }}</title>
  </head>
  <body>
    <h1>{{ .Item.Title }}</h1>
    <h2>{{ .Item.URL }}</h2>

    <div class="content">
      {{ .Content }}
    </div>
  </body>
</html>
`

type ReadCtx struct {
	Title   string
	Item    data.Meta
	Content template.HTML
}
