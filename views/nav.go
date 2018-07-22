package views

const nav = `<header>
  <h1><a href="/">papiermache</a></h1>
  <nav>
    <a href="/" {{ if eq . "toread" }}class="selected"{{ end }}>To Read</a>
    <a href="/liked" {{ if eq . "liked" }}class="selected"{{ end }}>Liked</a>
    <a href="/archived" {{ if eq . "archived" }}class="selected"{{ end }}>Archived</a>
    <span>|</span>
    <a href="">Bookmarklet</a>
    <a href="/generate" download>Download</a>
  </nav>
</header>`
