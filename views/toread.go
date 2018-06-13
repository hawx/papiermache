package views

import "hawx.me/code/papiermache/data"

const toRead = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    {{ template "nav" "toread" }}
    <ul>
      {{ range .Items }}
      <li>
        <a href="/read/{{ .Id }}">{{ .URL }}</a>
        {{ template "actions" . }}
      </li>
      {{ end }}
    </ul>
  </body>
</html>`

type ToReadCtx struct {
	Items []data.Meta
}
