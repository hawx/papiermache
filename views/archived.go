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
        {{ template "actions" . }}
        <a href="/read/{{ .Id }}">{{ .Title }}</a>
      </li>
      {{ end }}
    </ul>
  </body>
</html>`

type ArchivedCtx struct {
	Items []data.Meta
}
