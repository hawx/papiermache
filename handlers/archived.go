package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
)

func Archived(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListArchived()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Archived.Execute(w, views.ArchivedCtx{
			Items: items,
		})
	}
}
