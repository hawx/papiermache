package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/data"
	"hawx.me/code/route"
)

func Like(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]
		on := r.FormValue("un") == ""

		if err := db.Like(id, on); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
