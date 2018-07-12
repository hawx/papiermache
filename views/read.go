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
    {{ template "nav" "" }}

    <article>
      <header>
        <h1>{{ .Item.Title }}</h1>
        <a href="{{ .Item.URL }}">{{ .Item.URL }}</a>
        <time>{{ .Item.Added }}</time>
        {{ template "actions" .Item }}
      </header>

      <div class="content">
        {{ .Content }}
      </div>
    </article>
  </body>
</html>
`

type ReadCtx struct {
	Title   string
	Item    data.Meta
	Content template.HTML
}
