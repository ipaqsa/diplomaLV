package db

import (
	"agent/pkg"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

func InitDB(path string) error {
	if !exists(path) && path[len(path)-2:] != "db" {
		return errors.New("database`s not found")
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil
	}
	defer db.Close()
	err = initBucket(db)
	if err != nil {
		return err
	}
	return nil
}

func initBucket(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Keys"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func openDB(path string) (*bolt.DB, error) {
	if !exists(path) && path[len(path)-2:] != "db" {
		return nil, errors.New("database`s not found")
	}
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertKeys(data []InsertData) error {
	ref, err := openDB(pkg.Config.DBPath)
	if err != nil {
		return err
	}
	defer ref.Close()
	err = ref.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keys"))
		for _, d := range data {
			err = b.Put([]byte(d.Key), []byte(d.Value))
			if err != nil {
				err = tx.Rollback()
				if err != nil {
					return err
				}
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func RemoveKey(key string) error {
	ref, err := openDB(pkg.Config.DBPath)
	if err != nil {
		return err
	}
	defer ref.Close()
	err = ref.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keys"))
		err := b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetKeys() (map[string]string, error) {
	ref, err := openDB(pkg.Config.DBPath)
	if err != nil {
		return nil, err
	}
	defer ref.Close()
	var data = make(map[string]string)
	err = ref.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keys"))
		err := b.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
