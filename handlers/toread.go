package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
)

func ToRead(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListToRead()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.ToRead.Execute(w, views.ListCtx{
			Items: items,
		})
	}
}
