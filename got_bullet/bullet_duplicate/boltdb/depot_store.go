package boltdb

import (
	"encoding/binary"
	"fmt"

	"go.etcd.io/bbolt"
	"vixac.com/got/bullet_duplicate/model"
)

func (b *BoltStore) DepotPut(appID int32, key int64, value string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucketName := []byte(fmt.Sprintf("pigeon:app:%d", appID))
		bkt, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return bkt.Put(encodeInt64(key), []byte(value))
	})
}

func (b *BoltStore) DepotGet(appID int32, key int64) (string, error) {
	var val []byte
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucketName := []byte(fmt.Sprintf("pigeon:app:%d", appID))
		bkt := tx.Bucket(bucketName)
		if bkt == nil {
			return fmt.Errorf("not found")
		}
		val = bkt.Get(encodeInt64(key))
		if val == nil {
			return fmt.Errorf("not found")
		}
		return nil
	})
	return string(val), err
}

func (b *BoltStore) DepotDelete(appID int32, key int64) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucketName := []byte(fmt.Sprintf("pigeon:app:%d", appID))
		bkt := tx.Bucket(bucketName)
		if bkt == nil {
			return nil
		}
		return bkt.Delete(encodeInt64(key))
	})
}

func (b *BoltStore) DepotPutMany(appID int32, items []model.DepotKeyValueItem) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucketName := []byte(fmt.Sprintf("pigeon:app:%d", appID))
		bkt, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		for _, item := range items {
			if err := bkt.Put(encodeInt64(item.Key), []byte(item.Value)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *BoltStore) DepotGetMany(appID int32, keys []int64) (map[int64]string, []int64, error) {
	results := make(map[int64]string)
	var missing []int64

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucketName := []byte(fmt.Sprintf("pigeon:app:%d", appID))
		bkt := tx.Bucket(bucketName)
		if bkt == nil {
			missing = keys
			return nil
		}
		for _, k := range keys {
			val := bkt.Get(encodeInt64(k))
			if val == nil {
				missing = append(missing, k)
			} else {
				results[k] = string(val)
			}
		}
		return nil
	})

	return results, missing, err
}

func encodeInt64(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
