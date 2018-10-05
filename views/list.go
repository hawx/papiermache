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
        <div>
          <a class="read" href="/read/{{ .Id }}">{{ .Title }}</a>
          <div class="meta">
            <a href="{{ .URL }}">{{ domain .URL }}</a>
            <span>•</span>
            <time>{{ humanDate .Added }}</time>
          </div>
        </div>
        {{ template "actions" . }}
       </li>
      {{ end }}
    </ul>
  </body>
</html>`
}

type ListCtx struct {
	Items []data.Meta
}
