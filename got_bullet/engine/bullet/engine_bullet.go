package bullet

import (
	"errors"

	"vixac.com/got/engine"
)

type EngineBullet struct {
}

func (e *EngineBullet) Summary(gid engine.GidLookup) (*engine.GotSummary, error) {
	return nil, errors.New("not implemeneted")

}
