package boltdb

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/satty9/pocket-tg-bot-go/pkg/repository"
)

type BoltLocal struct {
	db *bolt.DB
}

// Refactoring. Rename this function. Maybe add it to Interface TokenRepositorier
func NewBoltDB(newDB *bolt.DB) *repository.TokenRepositorier {
	var tokenRepository repository.TokenRepositorier = &BoltLocal{db: newDB}
	return &tokenRepository
}

// implements interface TokenRepositorier
func (r *BoltLocal) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(IntToBytes(chatID), []byte(token))
	})
}

// implements interface TokenRepositorier
func (r *BoltLocal) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(IntToBytes(chatID))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func IntToBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
