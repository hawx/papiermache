package data

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

var (
	// The structure is pretty simple, I think.

	// These three buckets store the data keyed on a uuid
	metaBucketName    = []byte("meta")
	contentBucketName = []byte("content")
	rawBucketName     = []byte("raw")

	// And these buckets store lists of ids with key and value of the id
	toReadBucketName   = []byte("toRead")
	likedBucketName    = []byte("liked")
	archivedBucketName = []byte("archive")
)

// Item ids are a combination of the current UNIX time and a UUID, this is
// overkill, but whatever.
func makeId() (id []byte, now time.Time, err error) {
	buf := new(bytes.Buffer)
	now = time.Now().UTC()

	if err = binary.Write(buf, binary.LittleEndian, now.Unix()); err != nil {
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, uuid.New()); err != nil {
		return
	}

	id = buf.Bytes()
	return
}

type Meta struct {
	Id    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`

	Added    time.Time `json:"added"`
	Liked    time.Time `json:"liked"`
	Archived time.Time `json:"archived"`
}

func (m Meta) IsLiked() bool {
	return m.Liked != time.Time{}
}

func (m Meta) IsArchived() bool {
	return m.Archived != time.Time{}
}

type Database interface {
	ToRead(meta Meta, content, raw string) (id string, err error)
	Like(id string) error
	Archive(id string) error

	ListToRead() ([]Meta, error)
	ListLiked() ([]Meta, error)
	ListArchived() ([]Meta, error)

	Get(id string) (Meta, string, error)

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
			archivedBucketName,
			metaBucketName,
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

func (d *database) ToRead(meta Meta, content, raw string) (id string, err error) {
	key, now, err := makeId()
	if err != nil {
		return
	}

	meta.Id = base64.URLEncoding.EncodeToString(key)
	meta.Added = now

	return meta.Id, d.db.Update(func(tx *bolt.Tx) error {
		var (
			toReadBucket  = tx.Bucket(toReadBucketName)
			metaBucket    = tx.Bucket(metaBucketName)
			contentBucket = tx.Bucket(contentBucketName)
			rawBucket     = tx.Bucket(rawBucketName)
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

		if err = metaBucket.Put(key, value); err != nil {
			return err
		}

		return toReadBucket.Put(key, key)
	})
}

func (d *database) Like(id string) error {
	key, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		var (
			likedBucket = tx.Bucket(likedBucketName)
			metaBucket  = tx.Bucket(metaBucketName)
		)

		if err := updateMeta(key, metaBucket, func(meta Meta) Meta {
			meta.Liked = time.Now().UTC()
			return meta
		}); err != nil {
			return err
		}

		return likedBucket.Put(key, key)
	})
}

func (d *database) Archive(id string) error {
	key, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		var (
			toReadBucket   = tx.Bucket(toReadBucketName)
			archivedBucket = tx.Bucket(archivedBucketName)
			metaBucket     = tx.Bucket(metaBucketName)
		)

		if err := updateMeta(key, metaBucket, func(meta Meta) Meta {
			meta.Archived = time.Now().UTC()
			return meta
		}); err != nil {
			return err
		}

		if err := toReadBucket.Delete(key); err != nil {
			return err
		}

		return archivedBucket.Put(key, key)
	})
}

func (d *database) ListToRead() ([]Meta, error) {
	return d.list(toReadBucketName)
}

func (d *database) ListLiked() ([]Meta, error) {
	return d.list(likedBucketName)
}

func (d *database) ListArchived() ([]Meta, error) {
	return d.list(archivedBucketName)
}

func (d *database) list(bucket []byte) (metas []Meta, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		var (
			listBucket = tx.Bucket(bucket)
			metaBucket = tx.Bucket(metaBucketName)
		)

		c := listBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var meta Meta
			json.Unmarshal(metaBucket.Get(v), &meta)
			metas = append(metas, meta)
		}

		return nil
	})

	return
}

func (d *database) Get(id string) (meta Meta, content string, err error) {
	key, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return
	}

	err = d.db.View(func(tx *bolt.Tx) error {
		var (
			metaBucket    = tx.Bucket(metaBucketName)
			contentBucket = tx.Bucket(contentBucketName)
		)

		content = string(contentBucket.Get(key))

		value := metaBucket.Get(key)
		if value == nil {
			return errors.New("what, that doesn't even exist")
		}

		return json.Unmarshal(value, &meta)
	})

	return
}

func (d *database) Close() error {
	return d.db.Close()
}

func updateMeta(key []byte, metaBucket *bolt.Bucket, f func(Meta) Meta) error {
	value := metaBucket.Get(key)
	if value == nil {
		return errors.New("what, that doesn't even exist")
	}

	var meta Meta
	if err := json.Unmarshal(value, &meta); err != nil {
		return err
	}

	value, err := json.Marshal(f(meta))
	if err != nil {
		return err
	}

	return metaBucket.Put(key, value)
}
