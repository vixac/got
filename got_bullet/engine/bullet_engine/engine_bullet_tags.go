package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

// VX:TODO test
func (e *EngineBullet) TagItem(lookup engine.GidLookup, tagLookup engine.TagLookup) error {
	//VX:TODO tag string or tag lookup?

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}

	summary := engine.SummaryId(gid.IntValue)
	res, err := e.SummaryStore.Fetch([]engine.SummaryId{summary})
	if err != nil || res == nil {
		return err
	}
	val, ok := res[summary]
	if !ok {
		return errors.New("missing summary")
	}

	//VX:Note for now we are treating all taglookups as literals. But eventually the user might be able to
	//pass in a tag id which fetches the tag content

	tagLiteralString := tagLookup.Input
	tagLiteral := engine.TagLiteral{
		Display: tagLiteralString,
		Token:   "", //not fucking with token yet.
	}

	tag := engine.Tag{
		Literal:    &tagLiteral,
		Identifier: nil,
	}

	//ignore if duplicate
	for _, t := range val.Tags {
		if t.EqualTo(tag) {
			return nil
		}
	}

	val.Tags = append(val.Tags, tag)
	return e.SummaryStore.UpsertSummary(summary, val)

}

// VX:TODO test
func (e *EngineBullet) RemoveTag(tagLookup engine.TagLookup, summary engine.SummaryId) error {
	res, err := e.SummaryStore.Fetch([]engine.SummaryId{summary})
	if err != nil || res == nil {
		return err
	}
	val, ok := res[summary]
	if !ok {
		return errors.New("missing summary")
	}

	//VX:TODO again we're treating taglookup as a literal only.
	tagLiteralString := tagLookup.Input
	tagLiteral := engine.TagLiteral{
		Display: tagLiteralString,
		Token:   "", //not fucking with token yet.
	}

	tag := engine.Tag{
		Literal:    &tagLiteral,
		Identifier: nil,
	}

	var newTags []engine.Tag
	//ignore if duplicate
	for _, t := range val.Tags {
		if !t.EqualTo(tag) {
			newTags = append(newTags, t)
		}
	}
	val.Tags = newTags
	return e.SummaryStore.UpsertSummary(summary, val)
}
