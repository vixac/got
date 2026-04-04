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
	"vixac.com/got/engine/grove_engine"
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
	path := os.Getenv("GOT_BOLT")
	if path == "" {
		log.Fatal("missing env GOT_BOLT, which should be the path to the got bolt file")
	}

	sqlitePath := os.Getenv("GOT_SQLITE")
	if sqlitePath == "" {
		log.Fatal("missing env GOT_SQLITE, which should be the path to the got sqlite file")
	}
	sqlite, err := sqlite_store.NewSQLiteStore(sqlitePath)
	if err != nil {
		log.Fatal(err)
	}

	localBullet := local_bullet.LocalBullet{
		Space: space,
		Store: sqlite,
	}

	ene, err := grove_engine.NewGroveEngine(&localBullet)
	//ene, err := bullet_engine.NewEngineBullet(&localBullet)

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
