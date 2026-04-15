package main

import (
	"fmt"
	"log"
	"os"

	sqlite_store "github.com/vixac/bullet/store/sqlite"
	"github.com/vixac/bullet/store/store_interface"

	"github.com/vixac/firbolg_clients/bullet/local_bullet"
	"github.com/vixac/firbolg_clients/bullet/rest_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	"vixac.com/got/engine/grove_engine"
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

	fmt.Printf("VX: Not using the local bullet %s\n", localBullet.Space)
	logger := log.New(os.Stdout, "", log.LstdFlags)
	option := rest_bullet.WithLogger(logger)
	restClient := rest_bullet.NewRestClient("http://localhost:80", space, option)

	fmt.Printf("VX: rest client %s\n", restClient.AppId)
	ene, err := grove_engine.NewGroveEngine(&localBullet)
	//ene, err := grove_engine.NewGroveEngine(restClient)

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
