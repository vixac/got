package grove_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (g *GroveEngine) MarkResolved(lookup []engine.GidLookup) error {
	return errors.New(" MarkResolved Not impl")
}
func (g *GroveEngine) EditTitle(lookup engine.GidLookup, newHeading string) error {
	return errors.New("EditTitle Not impl")
}

func (g *GroveEngine) ScheduleItem(lookup engine.GidLookup, dateLookup engine.DateLookup) error {
	return errors.New("ScheduleItem Not impl")
}
func (g *GroveEngine) TagItem(lookup engine.GidLookup, tag engine.TagLookup) error {
	return errors.New("TagItem Not impl")
}
func (g *GroveEngine) ToggleCollapse(lookup engine.GidLookup, collapsed bool) error {
	return errors.New("ToggleCollapse Not impl")
}
