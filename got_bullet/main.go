package main

import (
	"vixac.com/got/cmd"
	"vixac.com/got/console"
	"vixac.com/got/engine/bullet"
)

func main() {
	engine := bullet.EngineBullet{}
	printer := console.Printer{}
	deps := cmd.RootDependencies{
		Printer: printer,
		Engine:  &engine,
	}
	cmd.Execute(deps)
	//println("VX: Hello got from go")
}
