package views

import "hawx.me/code/papiermache/data"

const liked = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    {{ template "nav" "liked" }}
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

type LikedCtx struct {
	Items []data.Meta
}
