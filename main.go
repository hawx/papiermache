package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/mauidude/go-readability"
	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/views"
	"hawx.me/code/route"
	"hawx.me/code/serve"
	"hawx.me/code/uberich"
)

func main() {
	var (
		settingsPath = flag.String("settings", "./settings.toml", "")
		port         = flag.String("port", "8080", "")
		socket       = flag.String("socket", "", "")
	)
	flag.Parse()

	var conf struct {
		Secret  string
		DbPath  string `toml:"database"`
		Uberich struct {
			AppName    string
			AppURL     string
			UberichURL string
			Secret     string
		}
	}
	if _, err := toml.DecodeFile(*settingsPath, &conf); err != nil {
		log.Fatal("toml:", err)
	}

	store := uberich.NewStore(conf.Secret)
	uberich := uberich.NewClient(conf.Uberich.AppName, conf.Uberich.AppURL, conf.Uberich.UberichURL, conf.Uberich.Secret, store)

	db, err := data.Open(conf.DbPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	shield := func(h http.HandlerFunc) http.Handler {
		return uberich.Protect(h, http.NotFoundHandler())
	}

	route.Handle("/", uberich.Protect(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			items, err := db.ListToRead()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")

			views.ToRead.Execute(w, views.ToReadCtx{
				Items: items,
			})
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")

			views.SignIn.Execute(w, nil)
		}),
	))

	route.Handle("/liked", shield(func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListLiked()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Liked.Execute(w, views.LikedCtx{
			Items: items,
		})
	}))

	route.Handle("/archived", shield(func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ListArchived()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		views.Archived.Execute(w, views.ArchivedCtx{
			Items: items,
		})
	}))

	route.Handle("/read/:id", shield(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	route.Handle("/add", shield(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	route.Handle("/like/:id", shield(func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		if err := db.Like(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}))

	route.Handle("/archive/:id", shield(func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		if err := db.Archive(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}))

	route.Handle("/sign-in", uberich.SignIn("/"))
	route.Handle("/sign-out", uberich.SignOut("/"))

	route.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, views.Styles)
	})

	serve.Serve(*port, *socket, route.Default)
}
