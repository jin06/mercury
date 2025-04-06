package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// your code here
}
