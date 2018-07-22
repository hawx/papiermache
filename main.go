package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/handlers"
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
		log.Println("data:", err)
		return
	}
	defer db.Close()

	shield := func(h http.HandlerFunc) http.Handler {
		return uberich.Protect(h, http.NotFoundHandler())
	}

	route.Handle("/", uberich.Protect(handlers.ToRead(db), handlers.SignIn()))
	route.Handle("/liked", shield(handlers.Liked(db)))
	route.Handle("/archived", shield(handlers.Archived(db)))
	route.Handle("/read/:id", shield(handlers.Read(db)))
	route.Handle("/add", shield(handlers.Add(db)))
	route.Handle("/like/:id", shield(handlers.Like(db)))
	route.Handle("/archive/:id", shield(handlers.Archive(db)))
	route.Handle("/generate", shield(handlers.Generate(db)))

	route.Handle("/sign-in", uberich.SignIn("/"))
	route.Handle("/sign-out", uberich.SignOut("/"))

	route.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, views.Styles)
	})

	serve.Serve(*port, *socket, route.Default)
}
