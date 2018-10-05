package views

const nav = `<header>
  <h1><a href="/">papiermache</a></h1>
  <nav>
    <div class="left">
      <a href="/" {{ if eq . "toread" }}class="selected"{{ end }}>To Read</a>
      <a href="/liked" {{ if eq . "liked" }}class="selected"{{ end }}>Liked</a>
      <a href="/archived" {{ if eq . "archived" }}class="selected"{{ end }}>Archived</a>
    </div>
    <div class="right">
      <a href="javascript:location.href='{{ baseURL }}/add?url='+encodeURIComponent(location.href)+'&redirect=origin;'">Bookmarklet</a>
      <a href="/generate" download>Download</a>
    </div>
  </nav>
</header>`
