package db

import (
	"github.com/boltdb/bolt"
)

const bucketName = "urls"
const bucketOffset = 1024 * 1024

// ServerDb holds current database connection
type ServerDb struct {
	db *bolt.DB
}

// Database is an generic interface for db operations
type Database interface {
	Read(key uint64) string
	Write(value string) (uint64, error)
}
