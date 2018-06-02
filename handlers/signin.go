package handlers

import (
	"net/http"

	"hawx.me/code/papiermache/views"
)

func SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		views.SignIn.Execute(w, nil)
	}
}
