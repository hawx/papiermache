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
        {{ template "actions" .Item }}
        <h1>{{ .Item.Title }}</h1>
        <div class="meta">
          <a href="{{ .Item.URL }}">{{ domain .Item.URL }}</a>
          <span>•</span>
          <time>{{ humanDate .Item.Added }}</time>
        </div>
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
