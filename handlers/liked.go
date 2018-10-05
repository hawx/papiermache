package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
)

func Liked(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListLiked()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Liked.Execute(w, views.ListCtx{
			Items: items,
		})
	}
}
