package data

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

var (
	toReadBucketName  = []byte("toRead")
	likedBucketName   = []byte("liked")
	archiveBucketName = []byte("archive")
	contentBucketName = []byte("content")
	rawBucketName     = []byte("raw")
)

type Meta struct {
	Id    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Database interface {
	ToRead(meta Meta, content, raw string) error
	Like(id string) error
	Archive(id string) error

	ListToRead() ([]Meta, error)
	ListLiked() ([]Meta, error)
	ListArchived() ([]Meta, error)

	GetToReadContent(id string) (Meta, string, error)

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
		for _, bucket := range [][]byte{
			toReadBucketName,
			likedBucketName,
			archiveBucketName,
			contentBucketName,
			rawBucketName,
		} {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return &database{db}, err
}

func (d *database) ToRead(meta Meta, content, raw string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		var (
			toReadBucket  = tx.Bucket(toReadBucketName)
			contentBucket = tx.Bucket(contentBucketName)
			rawBucket     = tx.Bucket(rawBucketName)
			key           = []byte(meta.Id)
		)

		value, err := json.Marshal(meta)
		if err != nil {
			return err
		}

		if err = contentBucket.Put(key, []byte(content)); err != nil {
			return err
		}

		if err = rawBucket.Put(key, []byte(raw)); err != nil {
			return err
		}

		return toReadBucket.Put(key, value)
	})
}

func (d *database) Like(id string) error {
	return d.move(id, toReadBucketName, likedBucketName)
}

func (d *database) Archive(id string) error {
	return d.move(id, toReadBucketName, archiveBucketName)
}

func (d *database) move(id string, from, to []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		var (
			fromBucket = tx.Bucket(from)
			toBucket   = tx.Bucket(to)
			key        = []byte(id)
		)

		value := fromBucket.Get(key)
		if err := fromBucket.Delete(key); err != nil {
			return err
		}

		return toBucket.Put(key, value)
	})
}

func (d *database) ListToRead() ([]Meta, error) {
	return d.list(toReadBucketName)
}

func (d *database) ListLiked() ([]Meta, error) {
	return d.list(toReadBucketName)
}

func (d *database) ListArchived() ([]Meta, error) {
	return d.list(toReadBucketName)
}

func (d *database) list(bucket []byte) (meta []Meta, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item Meta
			json.Unmarshal(v, &item)
			meta = append(meta, item)
		}

		return nil
	})

	return
}

func (d *database) GetToReadContent(id string) (Meta, string, error) {
	return d.getContent(id, toReadBucketName)
}

func (d *database) getContent(id string, bucket []byte) (item Meta, content string, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		var (
			metaBucket    = tx.Bucket(bucket)
			contentBucket = tx.Bucket(contentBucketName)
			key           = []byte(id)
		)

		content = string(contentBucket.Get(key))

		value := metaBucket.Get(key)
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
