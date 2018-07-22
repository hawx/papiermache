package views

import "hawx.me/code/papiermache/data"

const archived = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    {{ template "nav" "archived" }}
    <ul>
      {{ range .Items }}
      <li>
        <a href="/read/{{ .Id }}">{{ .Title }}</a>
        {{ template "actions" . }}
      </li>
      {{ end }}
    </ul>
  </body>
</html>`

type ArchivedCtx struct {
	Items []data.Meta
}
