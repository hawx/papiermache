package views

import (
	"html/template"

	"hawx.me/code/papiermache/data"
)

const generate = `<!DOCTYPE html>
  <body>
    <article>
      <header>
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
`

type GenerateCtx struct {
	Title   string
	Item    data.Meta
	Content template.HTML
}
