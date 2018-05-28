package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mauidude/go-readability"
	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
	"hawx.me/code/route"
	"hawx.me/code/serve"
)

func main() {
	var (
		dbPath = flag.String("db", "./db", "")
		port   = flag.String("port", "8080", "")
		socket = flag.String("socket", "", "")
	)
	flag.Parse()

	db, err := data.Open(*dbPath)
	if err != nil {
		log.Println(err)
		return
	}

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListToRead()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.ToRead.Execute(w, views.ToReadCtx{
			Items: items,
		})
	})

	route.HandleFunc("/liked", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListLiked()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Liked.Execute(w, views.LikedCtx{
			Items: items,
		})
	})

	route.HandleFunc("/archived", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListArchived()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Archived.Execute(w, views.ArchivedCtx{
			Items: items,
		})
	})

	route.HandleFunc("/read/:id", func(w http.ResponseWriter, r *http.Request) {
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
	})

	route.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		itemURL := r.FormValue("url")

		resp, err := http.Get(itemURL)
		if err != nil {
			http.Error(w, "Could not get '"+itemURL+"'", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Could not get '"+itemURL+"'", http.StatusBadRequest)
			return
		}

		doc, err := readability.NewDocument(string(html))
		if err != nil {
			http.Error(w, "Could not understand '"+itemURL+"'", http.StatusBadRequest)
			return
		}

		id, err := db.ToRead(data.Meta{
			URL:   itemURL,
			Title: "",
		}, doc.Content(), string(html))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/read/"+id, http.StatusFound)
	})

	route.HandleFunc("/like/:id", func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		if err := db.Like(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	route.HandleFunc("/archive/:id", func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		if err := db.Archive(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	route.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, views.Styles)
	})

	serve.Serve(*port, *socket, route.Default)
}
