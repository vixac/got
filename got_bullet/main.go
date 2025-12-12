package main

import (
	"log"
	"os"

	boldtb "github.com/vixac/bullet/store/boltdb"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	bullet_engine "vixac.com/got/engine/bullet_engine"
)

func main() {
	path := os.Getenv("GOT_BOLT")
	if path == "" {
		log.Fatal("missing env GOT_BOLT, which should be the path to the got bolt file")
	}
	boltdb, err := boldtb.NewBoltStore(path)
	if err != nil {
		log.Fatal(err)
	}
	localBullet := local_bullet.LocalBullet{
		AppId: 123,
		Store: boltdb,
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
