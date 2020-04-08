package db

import (
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vladikan/url-shortener/logger"
)

var db *bolt.DB
var mux sync.Mutex

const bucketName = "urls"

// Open will open bolt db connection
func Open() {
	blt, err := bolt.Open("urls.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Fatalf("Error while opening database, %s", err)
	}

	logger.Debug("Database opened")
	db = blt
}

// Close will close bolt db connection
func Close() {
	err := db.Close()
	if err != nil {
		logger.Fatalf("Error while closing database, %s")
	}
}

// Read will read stored value by its key
func Read(key string) string {
	var value string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(key))

		if v != nil && len(v) > 0 {
			value = string(v)
		}

		return nil
	})

	return value
}

// Write will store new key-value pair
func Write(key string, value string) error {
	mux.Lock()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})

	mux.Unlock()

	return err
}
