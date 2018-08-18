package views

import (
	"html/template"
	"io"
	"net/url"
	"time"
)

var SignIn, ToRead, Liked, Archived, Read interface {
	Execute(io.Writer, interface{}) error
}

var BaseURL string = ""

func init() {
	var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"baseURL": func() string { return BaseURL },
		"humanDate": func(t time.Time) string {
			return t.Format("January 02, 2006")
		},
		"domain": func(u string) string {
			pu, err := url.Parse(u)
			if err != nil {
				return u
			}
			return pu.Host
		},
	}).Parse(""))
	tmpl = template.Must(tmpl.New("toRead").Parse(toRead))
	tmpl = template.Must(tmpl.New("liked").Parse(liked))
	tmpl = template.Must(tmpl.New("archived").Parse(archived))
	tmpl = template.Must(tmpl.New("read").Parse(read))
	tmpl = template.Must(tmpl.New("signIn").Parse(signIn))

	tmpl = template.Must(tmpl.New("head").Parse(head))
	tmpl = template.Must(tmpl.New("actions").Parse(actions))
	tmpl = template.Must(tmpl.New("nav").Parse(nav))

	ToRead = &wrappedTemplate{tmpl, "toRead"}
	Liked = &wrappedTemplate{tmpl, "liked"}
	Archived = &wrappedTemplate{tmpl, "archived"}
	Read = &wrappedTemplate{tmpl, "read"}
	SignIn = &wrappedTemplate{tmpl, "signIn"}
}

type wrappedTemplate struct {
	t *template.Template
	n string
}

func (w *wrappedTemplate) Execute(wr io.Writer, data interface{}) error {
	return w.t.ExecuteTemplate(wr, w.n, data)
}
