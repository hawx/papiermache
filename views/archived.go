package views

import "hawx.me/code/papiermache/data"

const archived = `<!DOCTYPE html>
<html>
  <head>
    {{ template "head" }}
  </head>
  <body>
    <header>
      <h1>papiermache</h1>
      <nav>
        <a href="/">To Read</a>
        <a href="/liked">Liked</a>
        <a href="/archived" class="selected">Archived</a>
      </nav>
    </header>
    <ul>
      {{ range .Items }}
      <li>
        <a href="/read/{{ .Id }}">{{ .URL }}</a>
        <div class="actions">
          <a href="/like/{{ .Id }}">Like</a>
        </div>
      </li>
      {{ end }}
    </ul>
  </body>
</html>`

type ArchivedCtx struct {
	Items []data.Meta
}
