package views

import (
	"html/template"
	"io"
)

var List, Read interface {
	Execute(io.Writer, interface{}) error
}

func init() {
	var tmpl = template.Must(template.New("list").Funcs(template.FuncMap{}).Parse(list))
	tmpl = template.Must(tmpl.New("read").Parse(read))

	List = &wrappedTemplate{tmpl, "list"}
	Read = &wrappedTemplate{tmpl, "read"}
}

type wrappedTemplate struct {
	t *template.Template
	n string
}

func (w *wrappedTemplate) Execute(wr io.Writer, data interface{}) error {
	return w.t.ExecuteTemplate(wr, w.n, data)
}
