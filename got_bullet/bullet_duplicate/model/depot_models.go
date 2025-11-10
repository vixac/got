package model

type DepotRequest struct {
	Key   int64  `json:"key,string"`
	Value string `json:"value,string"`
}

type DepotKeyValueItem struct {
	Key   int64  `json:"key"`
	Value string `json:"value"`
}

type DepotPutManyRequest struct {
	Items []DepotKeyValueItem `json:"items"`
}

type DepotGetManyRequest struct {
	Keys []string `json:"keys"`
}

type DepotGetManyResponse struct {
	Values  map[int64]string `json:"values"`
	Missing []int64          `json:"missing"`
}
