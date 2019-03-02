package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"hawx.me/code/indieauth"
	"hawx.me/code/indieauth/sessions"
	"hawx.me/code/papiermache/data"
	"hawx.me/code/papiermache/handlers"
	"hawx.me/code/papiermache/views"
	"hawx.me/code/route"
	"hawx.me/code/serve"
)

func main() {
	var (
		settingsPath = flag.String("settings", "./settings.toml", "")
		port         = flag.String("port", "8080", "")
		socket       = flag.String("socket", "", "")
	)
	flag.Parse()

	var conf struct {
		BaseURL string
		Secret  string
		DbPath  string `toml:"database"`
		Me      string
	}
	if _, err := toml.DecodeFile(*settingsPath, &conf); err != nil {
		log.Fatal("toml:", err)
	}

	views.BaseURL = conf.BaseURL

	auth, err := indieauth.Authentication(conf.BaseURL, conf.BaseURL+"/callback")
	if err != nil {
		log.Fatal(err)
	}

	session, err := sessions.New(conf.Me, conf.Secret, auth)
	if err != nil {
		log.Fatal(err)
	}

	db, err := data.Open(conf.DbPath)
	if err != nil {
		log.Println("data:", err)
		return
	}
	defer db.Close()

	route.Handle("/", session.Choose(handlers.ToRead(db), handlers.SignIn()))
	route.Handle("/liked", session.Shield(handlers.Liked(db)))
	route.Handle("/archived", session.Shield(handlers.Archived(db)))
	route.Handle("/read/:id", session.Shield(handlers.Read(db)))
	route.Handle("/add", session.Shield(handlers.Add(db)))
	route.Handle("/like/:id", session.Shield(handlers.Like(db)))
	route.Handle("/archive/:id", session.Shield(handlers.Archive(db)))
	route.Handle("/generate", session.Shield(handlers.Generate(db)))

	route.Handle("/sign-in", session.SignIn())
	route.Handle("/callback", session.Callback())
	route.Handle("/sign-out", session.SignOut())

	route.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, views.Styles)
	})

	serve.Serve(*port, *socket, route.Default)
}
