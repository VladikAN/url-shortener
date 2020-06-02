package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vladikan/url-shortener/logger"
)

// Open will open bolt db connection
func Open() *ServerDb {
	blt, err := bolt.Open("urls.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Fatalf("Error while opening database, %s", err)
	}

	logger.Debug("Database opened")
	return &ServerDb{db: blt}
}

// Close will close bolt db connection
func Close(db Database) {
	blt, _ := db.(*ServerDb)
	err := blt.db.Close()
	if err != nil {
		logger.Fatalf("Error while closing database, %s")
	}
}

// Read will read stored value by its key
func (db ServerDb) Read(key uint64) string {
	var value string

	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil // first call nothing in store yet
		}

		v := b.Get(itob(key))

		if v != nil && len(v) > 0 {
			value = string(v)
		}

		return nil
	})

	return value
}

// Write will store new key-value pair
func (db ServerDb) Write(value string) (uint64, error) {
	db.mux.Lock()

	var id uint64
	err := db.db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
		id, _ = b.NextSequence()
		id += bucketOffset

		err := b.Put(itob(id), []byte(value))
		return err
	})

	db.mux.Unlock()

	return id, err
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
