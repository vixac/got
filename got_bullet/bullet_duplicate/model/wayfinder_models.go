package model

type WayFinderQueryItem struct {
	Key     string   `json:"key"`
	ItemId  int64    `json:"itemId"`
	Tag     *int64   `json:"tag,omitempty"`
	Metric  *float64 `json:"value,omitempty"`
	Payload string   `json:"payload"`
}

type WayFinderPutRequest struct {
	BucketId int32    `json:"bucketId"`
	Key      string   `json:"key"`
	Payload  string   `json:"payload"`
	Tag      *int64   `json:"tag,omitempty"`
	Metric   *float64 `json:"metric,omitempty"`
}

type WayFinderPrefixQueryRequest struct {
	BucketId   int32    `json:"bucketId"`
	Prefix     string   `json:"prefix"`
	Tags       []int64  `json:"tags,omitempty"`
	Metric     *float64 `json:"metric,omitempty"`
	MetricIsGt bool     `json:"metricIsGt"`
}

type WayFinderGetResponse struct {
	ItemId  int64    `json:"itemId"`
	Payload string   `json:"payload"`
	Tag     *int64   `json:"tag,omitempty"`
	Metric  *float64 `json:"metric,omitempty"`
}

type WayFinderGetOneRequest struct {
	BucketId int32  `json:"bucketId"`
	Key      string `json:"key"`
}
