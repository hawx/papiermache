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
        {{ template "actions" . }}
        <a href="/read/{{ .Id }}">{{ .Title }}</a>
      </li>
      {{ end }}
    </ul>
  </body>
</html>`

type LikedCtx struct {
	Items []data.Meta
}
