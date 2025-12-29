package main

import (
	"log"
	"os"

	boldtb "github.com/vixac/bullet/store/boltdb"
	"github.com/vixac/bullet/store/store_interface"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	bullet_engine "vixac.com/got/engine/bullet_engine"
)

func main() {
	path := os.Getenv("GOT_BOLT")
	//fmt.Printf("VX: LOOKING AT PATH %s\n", path)
	if path == "" {
		log.Fatal("missing env GOT_BOLT, which should be the path to the got bolt file")
	}
	bolt, err := boldtb.NewBoltStore(path)
	if err != nil {
		log.Fatal(err)
	}

	space := store_interface.TenancySpace{
		AppId:     123,
		TenancyId: 0,
	}

	/*
		bolt.MigrateToTenantBuckets(space, 1001)
		bolt.MigrateToTenantBuckets(space, 100)
		bolt.MigrateToTenantBuckets(space, 1002)
		bolt.MigrateToTenantBuckets(space, 1003)
		bolt.MigrateToTenantBuckets(space, 1004)
		bolt.MigrateToTenantBuckets(space, 0)
		bolt.MigrateToTenantBuckets(space, 1005)
	*/
	localBullet := local_bullet.LocalBullet{
		Space: space,
		Store: bolt,
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
