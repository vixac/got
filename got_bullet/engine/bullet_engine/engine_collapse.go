package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (e *EngineBullet) ToggleCollapse(lookup engine.GidLookup, collapse bool) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	summaryId := engine.SummaryId(gid.IntValue)
	summaries, err := e.SummaryStore.Fetch([]engine.SummaryId{summaryId})
	if err != nil {
		return err
	}
	summary, ok := summaries[summaryId]
	if !ok {
		return errors.New("missing summary, can't collapse it")
	}
	if summary.Flags == nil {
		summary.Flags = make(map[string]bool)
	}
	if collapse {
		summary.Flags["collapsed"] = true
	} else {
		summary.Flags["collapsed"] = false
	}

	err = e.SummaryStore.UpsertSummary(summaryId, summary)
	if err != nil {
		return err
	}
	return nil
}
