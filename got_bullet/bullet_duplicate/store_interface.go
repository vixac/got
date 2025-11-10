package bullet_duplicate

import (
	"vixac.com/got/bullet_duplicate/model"
)

// VX:TODO this has been copy pasted from bullet. The plan is to import all via a submodule provided by bullet.
type TrackStore interface {
	TrackPut(appID int32, bucketID int32, key string, value int64, tag *int64, metric *float64) error
	TrackGet(appID int32, bucketID int32, key string) (int64, error)
	TrackDelete(appID int32, bucketID int32, key string) error
	TrackClose() error
	TrackPutMany(appID int32, items map[int32][]model.TrackKeyValueItem) error
	TrackGetMany(appID int32, keys map[int32][]string) (map[int32]map[string]model.TrackValue, map[int32][]string, error)
	GetItemsByKeyPrefix(
		appID, bucketID int32,
		prefix string,
		tags []int64, // optional slice of tags
		metricValue *float64, // optional metric value
		metricIsGt bool, // "gt" or "lt"
	) ([]model.TrackKeyValueItem, error)
}

type DepotStore interface {
	DepotPut(appID int32, key int64, value string) error
	DepotGet(appID int32, key int64) (string, error)
	DepotDelete(appID int32, key int64) error
	DepotPutMany(appID int32, items []model.DepotKeyValueItem) error
	DepotGetMany(appID int32, keys []int64) (map[int64]string, []int64, error)
}

// using its own ids, wayfinder uses track and depot to provide a query to payload interface.
type WayFinderStore interface {
	WayFinderPut(appID int32, bucketID int32, key string, payload string, tag *int64, metric *float64) (int64, error)
	WayFinderGetByPrefix(appID int32, bucketID int32,
		prefix string,
		tags []int64, // optional slice of tags
		metricValue *float64, // optional metric value
		metricIsGt bool, // "gt" or "lt"
	) ([]model.WayFinderQueryItem, error)

	WayFinderGetOne(appID int32, bucketID int32, key string) (*model.WayFinderGetResponse, error)
}

type Store interface {
	TrackStore
	DepotStore
	WayFinderStore
}
