package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/route"
)

func Archive(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		if err := db.Archive(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
