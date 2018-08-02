package handlers

// https://stackoverflow.com/questions/5379565/kindle-periodical-format

import (
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"

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
	}
}

func generateMobi(db data.Database, w io.Writer) error {
	list, err := db.ListToRead()
	if err != nil {
		return err
	}

	var articles []periodize.Article

	encoder := encoding.ReplaceUnsupported(charmap.ISO8859_1.NewEncoder())

	for _, item := range list {
		meta, body, err := db.Get(item.Id)
		if err != nil {
			log.Println("generateMobi: could not find", item.Id)
			continue
		}

		encodedBody, err := encoder.String(body)
		if err != nil {
			log.Println("generateMobi: failed to encode body of", item.Id, ":", err)
			continue
		}

		articles = append(articles, periodize.Article{
			Title:   meta.Title,
			Author:  "",
			Content: "<body>" + encodedBody + "</body>",
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
