package views

import "hawx.me/code/papiermache/data"

const toRead = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    <header>
      <h1>papiermache</h1>
      <nav>
        <a href="/" class="selected">To Read</a>
        <a href="/liked">Liked</a>
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

type ToReadCtx struct {
	Items []data.Meta
}
