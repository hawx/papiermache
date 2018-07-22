package handlers

// https://stackoverflow.com/questions/5379565/kindle-periodical-format

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/djcrock/periodize"
	"github.com/google/uuid"
	"hawx.me/code/papiermache/data"
)

func Generate(db data.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=papiermache.mobi")
		w.Header().Set("Content-Type", "application/x-mobipocket-ebook")

		if err := generateMobi(db, w); err != nil {
			log.Println(err)
		}

		w.WriteHeader(200)
	}
}

func generateMobi(db data.Database, w io.Writer) error {
	list, err := db.ListToRead()
	if err != nil {
		return err
	}

	var articles []periodize.Article

	for _, item := range list {
		meta, body, err := db.Get(item.Id)
		if err != nil {
			continue
		}

		articles = append(articles, periodize.Article{
			Title:   meta.Title,
			Author:  "",
			Content: "<body>" + body + "</body>",
		})
	}

	issue := periodize.Issue{
		UniqueID:    uuid.New().String(),
		Title:       "My Periodical",
		Creator:     "papiermache",
		Publisher:   "papiermache",
		Subject:     "",
		Description: "Papiermache reading list",
		Date:        time.Now().Format("2006-01-02"),
		Sections: []periodize.Section{
			{
				Title:    "To Read",
				Articles: articles,
			},
		},
	}

	return issue.GenerateMobi(w)
}
