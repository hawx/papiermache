package views

import "hawx.me/code/papiermache/data"

const list = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>{{ .Title }}</title>
  </head>
  <body>
    <h1>{{ .Title }}</h1>
    <ul>
      {{ range .Items }}
      <li><a href="/read/{{ .Id }}">{{ .URL }}</a></li>
      {{ end }}
    </ul>
  </body>
</html>`

type ListCtx struct {
	Title string
	Items []data.Meta
}
