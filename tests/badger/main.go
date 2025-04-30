package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	// db, err := badger.Open(badger.DefaultOptions("./badger"))
	db, err := badger.Open(badger.DefaultOptions("/Users/jinlong/tmp/mercury-badger"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("answer"), []byte("42")).WithMeta(byte(1))
		err := txn.SetEntry(e)
		return err
	})
	if err != nil {
		log.Fatal()
	}

	defer db.Close()

	// your code here
}
