package main

import (
	"log"

	boldtb "github.com/vixac/bullet/store/boltdb"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	"vixac.com/got/engine/bullet"
)

func main() {
	//inter := bullet_interface.BulletClientInterface{}
	boltdb, err := boldtb.NewBoltStore("boom")
	if err != nil {
		log.Fatal(err)
	}
	localBullet := local_bullet.LocalBullet{
		AppId: 123,
		Store: boltdb,
	}
	ene := bullet.EngineBullet{
		Client: &localBullet,
	}
	printer := console.Printer{}
	deps := cmd.RootDependencies{
		Printer: printer,
		Engine:  &ene,
	}
	cmd.Execute(deps)
	//println("VX: Hello got from go")
}
