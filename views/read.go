package views

import (
	"html/template"

	"hawx.me/code/papiermache/data"
)

const read = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    <header>
      <h1>{{ .Item.Title }}</h1>
      <h2>{{ .Item.URL }}</h2>
      <time>{{ .Item.Added }}</time>
    </header>

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
