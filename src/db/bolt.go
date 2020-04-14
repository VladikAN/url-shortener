package db

import (
	"encoding/binary"
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
func Read(key uint64) string {
	var value string
	logger.Debugf("Request for %d", key)

	db.View(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
		v := b.Get(itob(key))

		if v != nil && len(v) > 0 {
			value = string(v)
		}

		return nil
	})

	return value
}

// Write will store new key-value pair
func Write(value string) (uint64, error) {
	mux.Lock()

	var id uint64
	err := db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
		id, _ = b.NextSequence()

		err := b.Put(itob(id), []byte(value))
		return err
	})

	mux.Unlock()

	return id, err
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
