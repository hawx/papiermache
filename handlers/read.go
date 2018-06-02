package handlers

import (
	"html/template"
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
	"hawx.me/code/route"
)

func Read(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		item, content, err := db.Get(id)
		if err != nil {
			// obviously do better than this
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		views.Read.Execute(w, views.ReadCtx{
			Title:   "Reading...",
			Item:    item,
			Content: template.HTML(content),
		})
	}
}
