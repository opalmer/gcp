package db

import (
	"github.com/boltdb/bolt"
	"log"
)

// Open - Opens the bolt database, checks for errors and
// returns bolt.DB
func Open(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
