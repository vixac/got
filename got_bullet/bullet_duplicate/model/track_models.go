package model

type TrackRequest struct {
	BucketID int32    `json:"bucketId"`
	Key      string   `json:"key"`
	Value    int64    `json:"value,string"`
	Tag      *int64   `json:"tag,omitempty"`
	Metric   *float64 `json:"metric,omitempty"`
}

type TrackPutManyRequest struct {
	Buckets []TrackPutItems `json:"buckets"`
}

type TrackGetManyRequest struct {
	Buckets []TrackGetKeys `json:"buckets"`
}

type MetricFilter struct {
	Operator string  `json:"operator"` // "gt", "lt", etc.
	Value    float64 `json:"value"`
}
type TrackGetItemsByPrefixRequest struct {
	BucketID int32         `json:"bucketId"`
	Prefix   string        `json:"prefix"`
	Tags     []int64       `json:"tags,omitempty"`   // optional IN clause
	Metric   *MetricFilter `json:"metric,omitempty"` // optional metric filter
}

type TrackKeyValueItem struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type TrackPutItems struct {
	BucketID int32               `json:"bucketId"`
	Items    []TrackKeyValueItem `json:"items"`
}

type TrackGetKeys struct {
	BucketID int32    `json:"bucketId"`
	Keys     []string `json:"keys"`
}

type TrackGetManyResponse struct {
	Values  map[string]map[string]TrackValue `json:"values"`  // bucketId -> (key -> value)
	Missing map[string][]string              `json:"missing"` // bucketId -> list of missing keys
}

type TrackValue struct {
	Value  int64    `bson:"value"`
	Tag    *int64   `bson:"tag,omitempty"`
	Metric *float64 `bson:"metric,omitempty"`
}
