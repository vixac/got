package grove_engine

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

const (
	groveIdGenBucket int32 = 200
	aliasBucket      int32 = 2001
	numberGoBucket   int32 = 2006
	longFormBucket   int32 = 2005
)

/*
	//GotAliasInterface
	GotCreateItemInterface
	GotFetchInterface
	RestoreInterface
	NoteInterface
	GotEditInterface
	GotTreeInterface
*/

type GroveEngine struct {
	Client        bullet_interface.BulletClientInterface
	AliasStore    engine_util.AliasStoreInterface
	GidLookup     engine.GidLookupInterface
	NumberGoStore engine.NumberGoStoreInterface
	LongFormStore engine.LongFormStoreInterface
	IdGenerator   engine.IdGeneratorInterface
}

func NewGroveEngine(client bullet_interface.BulletClientInterface) (*GroveEngine, error) {
	numberGoCodec := &engine_util.JSONCodec[engine_util.NumberGoBlock]{}
	numberGoStore, err := engine_util.NewBulletNumberGoStore(numberGoBucket, client, client, numberGoCodec)
	if err != nil {
		return nil, err
	}
	longFormStore, err := engine_util.NewBulletLongFormStore(longFormBucket, client, client)
	if err != nil {
		return nil, err
	}

	aliasStore, err := engine_util.NewBulletAliasStore(client, aliasBucket)
	if err != nil {
		return nil, err
	}
	idGenerator := engine_util.NewIdBulletGenerator(client, groveIdGenBucket, "next-id-list", "", "latest")

	gidLookup, err := engine_util.NewBulletGidLookup(aliasStore, numberGoStore, idGenerator)
	if err != nil {
		return nil, err
	}

	return &GroveEngine{
		Client:        client,
		GidLookup:     gidLookup,
		AliasStore:    aliasStore,
		NumberGoStore: numberGoStore,
		LongFormStore: longFormStore,
		IdGenerator:   idGenerator,
	}, nil
}
