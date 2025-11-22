package main

import (
	"log"

	boldtb "github.com/vixac/bullet/store/boltdb"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	bullet_engine "vixac.com/got/engine/bullet_engine"
)

func main() {
	boltdb, err := boldtb.NewBoltStore("got-bolt")
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
	//println("VX: Hello got from go")
}
