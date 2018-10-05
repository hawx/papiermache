package views

import "hawx.me/code/papiermache/data"

func list(of string) string {
	return `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    {{ template "nav" "` + of + `" }}
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
}

type ListCtx struct {
	Items []data.Meta
}
