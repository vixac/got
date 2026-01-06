package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

// VX:TODO test
func (e *EngineBullet) AddTag(tag string, summary engine.SummaryId) error {
	res, err := e.SummaryStore.Fetch([]engine.SummaryId{summary})
	if err != nil || res == nil {
		return err
	}
	val, ok := res[summary]
	if !ok {
		return errors.New("missing summary")
	}

	//ignore if duplicate
	for _, t := range val.Tags {
		if t == tag {
			return nil
		}
	}
	val.Tags = append(val.Tags, tag)
	return e.SummaryStore.UpsertSummary(summary, val)

}

func (e *EngineBullet) RemoveTag(tag string, summary engine.SummaryId) error {
	return nil
}
