package engine

import (
	"github.com/vixac/bullet/store/ram"
	"github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/local_bullet"
)

// Instead of a true mock client, we just create a real client using RAM storage.
func BuildTestClient() bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	space := store_interface.TenancySpace{
		AppId:     12,
		TenancyId: 100,
	}
	localClient := &local_bullet.LocalBullet{
		Store: store,
		Space: space,
	}
	return localClient
}
