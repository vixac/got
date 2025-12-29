package main

import (
	"log"
	"os"

	"github.com/vixac/bullet/store/boltdb"
	"github.com/vixac/bullet/store/store_interface"

	"github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	"vixac.com/got/engine/bullet_engine"
)

func main() {
	space := store_interface.TenancySpace{
		AppId:     123,
		TenancyId: 0,
	}
	path := os.Getenv("GOT_BOLT")
	if path == "" {
		log.Fatal("missing env GOT_BOLT, which should be the path to the got bolt file")
	}
	bolt, err := boltdb.NewBoltStore(path)
	if err != nil {
		log.Fatal(err)
	}

	/*
		sqlite, err := sqlite_store.NewSQLiteStore("got-sqlite.sqlite3")
		if err != nil {
			log.Fatal(err)
		}
	*/
	localBullet := local_bullet.LocalBullet{
		Space: space,
		Store: bolt, //sqlite,
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
