package pkg

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

var (
	Bucket = []byte("bucket")
	IndexKey = []byte("index")
	EmptyIndex uint64 = 0
)

type Storage struct {
	Path string
	DB   *bolt.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}

	storage := Storage{
		Path: path,
		DB:   db,
	}
	return &storage, nil
}

func (storage *Storage) Create(index uint64) error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(Bucket)
		if err != nil {
			return err
		}

		bs := bucket.Get(IndexKey)
		if len(bs) != 0 {
			return nil
		}
		bs, err = jsoniter.Marshal(index)
		if err != nil {
			return err
		}
		return bucket.Put(IndexKey, bs)
	})
}

func (storage *Storage) Put(index uint64) error {
	bs, err := jsoniter.Marshal(index)
	if err != nil {
		return err
	}
	return storage.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(Bucket)
		return bucket.Put(IndexKey, bs)
	})
}

func (storage *Storage) Get() (uint64, error) {
	var bs []byte
	err := storage.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(Bucket)
		bs = bucket.Get(IndexKey)
		return nil
	})
	if len(bs) == 0 || err != nil {
		return EmptyIndex, err
	}

	var index uint64
	err = jsoniter.Unmarshal(bs, &index)
	return index, err
}

func (storage *Storage) Stop() error {
	return storage.DB.Close()
}
