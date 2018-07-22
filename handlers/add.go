package handlers

import (
	"net/http"

	"github.com/antchfx/goreadly"
	"hawx.me/code/papiermache/data"
)

func Add(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemURL := r.FormValue("url")

		resp, err := http.Get(itemURL)
		if err != nil {
			http.Error(w, "Could not get '"+itemURL+"'", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		// html, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	http.Error(w, "Could not get '"+itemURL+"'", http.StatusBadRequest)
		// 	return
		// }

		doc, err := goreadly.ParseResponse(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// doc, err := readability.NewDocument(string(html))
		// if err != nil {
		// 	http.Error(w, "Could not understand '"+itemURL+"'", http.StatusBadRequest)
		// 	return
		// }

		id, err := db.ToRead(data.Meta{
			URL:   itemURL,
			Title: doc.Title,
		}, doc.Body, "no raw body")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/read/"+id, http.StatusFound)
	}
}
