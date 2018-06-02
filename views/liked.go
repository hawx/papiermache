package views

import "hawx.me/code/papiermache/data"

const liked = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    <header>
      <h1><a href="/">papiermache</a></h1>
      <nav>
        <a href="/">To Read</a>
        <a href="/liked" class="selected">Liked</a>
        <a href="/archived">Archived</a>
      </nav>
    </header>
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

type LikedCtx struct {
	Items []data.Meta
}
