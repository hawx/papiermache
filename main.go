package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/mauidude/go-readability"
	"hawx.me/code/route"
	"hawx.me/code/serve"
)

var itemBucketName = []byte("items")

type Item struct {
	Id      string `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
	HTML    string `json:"html"`
}

type Database interface {
	List() ([]Item, error)
	Put(Item) error
	Get(id string) (Item, error)
	Close() error
}

type database struct {
	db *bolt.DB
}

func Open(path string) (Database, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(itemBucketName)
		return err
	})

	return &database{db}, err
}

func (d *database) List() (items []Item, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(itemBucketName)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item Item
			json.Unmarshal(v, &item)
			items = append(items, item)
		}

		return nil
	})

	return
}

func (d *database) Put(item Item) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(itemBucketName)

		key := []byte(item.Id)
		value, err := json.Marshal(item)

		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
}

func (d *database) Get(id string) (item Item, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(itemBucketName)

		value := b.Get([]byte(id))
		if value == nil {
			return errors.New("what, that doesn't even exist")
		}

		return json.Unmarshal(value, &item)
	})

	return
}

func (d *database) Close() error {
	return d.db.Close()
}

func main() {
	var (
		dbPath = flag.String("db", "./db", "")
		port   = flag.String("port", "8080", "")
		socket = flag.String("socket", "", "")
	)
	flag.Parse()

	db, err := Open(*dbPath)
	if err != nil {
		log.Println(err)
		return
	}

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		fmt.Fprint(w, "<ul>")
		for _, item := range items {
			fmt.Fprintf(w, `<li>
  <a href="/read/%s">%s</a>
</li>`, item.Id, item.URL)
		}
		fmt.Fprint(w, "</ul>")
	})

	route.HandleFunc("/read/:id", func(w http.ResponseWriter, r *http.Request) {
		id := route.Vars(r)["id"]

		item, err := db.Get(id)
		if err != nil {
			// obviously do better than this
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, item.Content)
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

		id := uuid.New().String()
		db.Put(Item{
			Id:      id,
			URL:     itemURL,
			Title:   "",
			HTML:    string(html),
			Content: doc.Content(),
		})

		http.Redirect(w, r, "/read/"+id, http.StatusFound)
	})

	serve.Serve(*port, *socket, route.Default)
}
