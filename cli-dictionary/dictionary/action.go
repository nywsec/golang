package dictionary

import (
	"bytes"
	"encoding/gob"
	"sort"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
)

func (d *Dictionary) Add(word, definition string) error {

	entry := Entry{
		Word:       strings.Title(word),
		Definition: definition,
		CreateAt:   time.Now(),
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(entry)

	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(word), buffer.Bytes())

	})
}

func (d *Dictionary) Remove(word string) error {

	return d.db.Update(func(txn *badger.Txn) error {

		return txn.Delete([]byte(word))
	})
}

func (d *Dictionary) Get(word string) (Entry, error) {

	var entry Entry
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(word))
		if err != nil {
			return err
		}

		entry, err = getEntry(item)
		return err

	})
	return entry, err
}

//list retrieve any Dictionary's items .
//[]string is a alphabetically sorted array with word.
//[string]Entry is a map of the words and their definition.

func (d *Dictionary) List() ([]string, map[string]Entry, error) {

	entries := make(map[string]Entry)
	err := d.db.View(func(tnx *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := tnx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			entry, err := getEntry(item)
			if err != nil {
				return err
			}
			entries[entry.Word] = entry
		}
		return nil
	})
	return sortedKeys(entries), entries, err
}

func sortedKeys(entries map[string]Entry) []string {

	keys := make([]string, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys

}

func getEntry(item *badger.Item) (Entry, error) {
	var entry Entry
	var buffer bytes.Buffer
	err := item.Value(func(val []byte) error {
		_, err := buffer.Write(val)
		return err
	})

	dec := gob.NewDecoder(&buffer)
	err = dec.Decode(&entry)
	return entry, err
}
