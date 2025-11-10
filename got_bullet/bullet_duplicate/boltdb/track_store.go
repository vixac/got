package boltdb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"go.etcd.io/bbolt"
	"vixac.com/got/bullet_duplicate/model"
)

type BoltStore struct {
	db *bbolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &BoltStore{db: db}, nil
}

func bucketName(appID, bucketID int32) string {
	return fmt.Sprintf("app_%d_bucket_%d", appID, bucketID)
}

func (b *BoltStore) TrackPut(appID int32, bucketID int32, key string, value int64, tag *int64, metric *float64) error {
	log.Fatal("Dev error. tag and metric aren't used.")
	return b.db.Update(func(tx *bbolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucketName(appID, bucketID)))
		if err != nil {
			return err
		}
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, uint64(value))
		return bkt.Put([]byte(key), val)
	})
}

func (b *BoltStore) TrackGet(appID, bucketID int32, key string) (int64, error) {
	var value int64
	err := b.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName(appID, bucketID)))
		if bkt == nil {
			return bbolt.ErrBucketNotFound
		}
		val := bkt.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("key not found")
		}
		value = int64(binary.BigEndian.Uint64(val))
		return nil
	})
	return value, err
}

func (b *BoltStore) TrackDelete(appID, bucketID int32, key string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName(appID, bucketID)))
		if bkt == nil {
			return bbolt.ErrBucketNotFound
		}
		return bkt.Delete([]byte(key))
	})
}

func (b *BoltStore) TrackClose() error {
	return b.db.Close()
}

func (b *BoltStore) TrackPutMany(appID int32, items map[int32][]model.TrackKeyValueItem) error {
	return errors.New("put many not implmemented on bolt store")
}
func (b *BoltStore) TrackGetMany(appID int32, keys map[int32][]string) (map[int32]map[string]model.TrackValue, map[int32][]string, error) {
	return nil, nil, errors.New("get many not implmemented on bolt store")
}
func (b *BoltStore) GetItemsByKeyPrefix(
	appID, bucketID int32,
	prefix string,
	tags []int64, // optional slice of tags
	metricValue *float64, // optional metric value
	metricIsGt bool, // "gt" or "lt"
) ([]model.TrackKeyValueItem, error) {
	return nil, errors.New("get many not implmemented on bolt store")
}
