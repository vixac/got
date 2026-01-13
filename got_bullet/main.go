package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vixac/bullet/store/store_interface"

	migrator "github.com/vixac/bullet/store/migrator"
	sqlite_store "github.com/vixac/bullet/store/sqlite"
	"github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	"vixac.com/got/engine/bullet_engine"
)

func migrateBucket(bucket int32, m *migrator.TrackMigrator) {
	fmt.Printf("VX: MIGRATION on bucket %d starting...\n", bucket)
	err := m.Migrate(bucket)
	if err != nil {
		fmt.Printf("VX: error migrating bucket %d, %s", bucket, err)
		log.Fatal()
	}

	fmt.Printf("VX: MIGRATION on bucket %d Complete.\n", bucket)

}
func main() {
	space := store_interface.TenancySpace{
		AppId:     123,
		TenancyId: 0,
	}
	/*
		path := os.Getenv("GOT_BOLT")
		if path == "" {
			log.Fatal("missing env GOT_BOLT, which should be the path to the got bolt file")
		}

		bolt, err := boltdb.NewBoltStore(path)
		if err != nil {
			log.Fatal(err)
		}
	*/
	sqlitePath := os.Getenv("GOT_SQLITE")
	if sqlitePath == "" {
		log.Fatal("missing env GOT_SQLITE, which should be the path to the got sqlite file")
	}
	sqlite, err := sqlite_store.NewSQLiteStore(sqlitePath)
	if err != nil {
		log.Fatal(err)
	}

	/*
		trackMigrator := migrator.TrackMigrator{
			SourceTrack: bolt,
			TargetTrack: sqlite,
			Tenancy:     space,
		}
		trackMigrator.Migrate(100)
		trackMigrator.Migrate(1001)
		trackMigrator.Migrate(1003)
		depotMigrator := migrator.DepotMigrator{
			SourceDepot: bolt,
			TargetDepot: sqlite,
			Tenancy:     space,
		}

		allItems, err := bolt.DepotGetAll(depotMigrator.Tenancy)
		if err != nil {
			log.Fatal("fetching all keys didn't work.")
		}

		var allKeys []int64
		for k, _ := range allItems {
			allKeys = append(allKeys, k)
		}
		totalKeys := len(allItems)
		fmt.Printf("VX: depoing this many keys %d\n", totalKeys)
		depotMigrator.Migrate(allKeys)
		//ok so I need to get all the keys. The migrator should know how to do it.
		//or a depot debug class or somethign.
		//bolt.DepotGetMany(space)
		//depotMigrator.Migrate()
		fmt.Printf("VX: THIS IS MIGRATED")
	*/
	localBullet := local_bullet.LocalBullet{
		Space: space,
		Store: sqlite,
	}

	ene, err := bullet_engine.NewEngineBullet(&localBullet)
	if err != nil {
		log.Fatal(err)
	}
	printer := console.Printer{}
	deps := cmd.RootDependencies{
		Printer: printer,
		Engine:  ene,
	}
	cmd.Execute(deps)
}
