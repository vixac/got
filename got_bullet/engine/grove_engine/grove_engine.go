package grove_engine

import "github.com/vixac/firbolg_clients/bullet/bullet_interface"

type GroveEngine struct {
	Client bullet_interface.BulletClientInterface
}

func NewGroveEngine(client bullet_interface.BulletClientInterface) GroveEngine {

	return GroveEngine{
		Client: client,
	}
}

func (g *GroveEngine) Wohoo() {

}
