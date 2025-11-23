package bullet_engine

import "github.com/vixac/firbolg_clients/bullet/bullet_interface"

type AggStoreInterface interface {
	UpsertAggregate()
}

type BulletAggStore struct {
}

func NewBulletAggStore(client bullet_interface.DepotClientInterface, namespace int32) (AggStoreInterface, error) {
	return nil, nil
}
