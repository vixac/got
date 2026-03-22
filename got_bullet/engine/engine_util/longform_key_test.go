package engine_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"vixac.com/got/engine"
)

func TestLongFormKey(t *testing.T) {

	firstNoteId := engine.FirstNoteId()
	assert.Equal(t, firstNoteId.IntValue, int64(1000))
	assert.Equal(t, firstNoteId.AasciValue, "rs")

	gotId, err := engine.NewGotIdFromInt(10)
	assert.NoError(t, err)
	assert.Equal(t, gotId.IntValue, int32(10))
	assert.Equal(t, gotId.AasciValue, "a")

	millis := "1774179116040"
	date, err := engine.EpochMillisStringToDate(millis)
	assert.NoError(t, err)
	longFormKey := engine.LongFormKey{
		NoteId:      firstNoteId,
		GotId:       *gotId,
		CreatedTime: *date,
	}
	str := longFormKey.ToString()
	assert.Equal(t, str, "a:rs:1774179116040")

	longFormIdRebuilt, err := engine.NewLongFormKeyFromString(str)
	assert.NoError(t, err)
	assert.Equal(t, longFormIdRebuilt.GotId.AasciValue, "a")
	assert.Equal(t, longFormIdRebuilt.NoteId.AasciValue, "rs")
	assert.Equal(t, engine.TimeToMillisString(longFormIdRebuilt.CreatedTime), millis)

	nextMillis := "1774179116041"
	nextTime, err := engine.EpochMillisStringToDate(nextMillis)
	assert.NoError(t, err)
	nextKey := longFormKey.Next(*nextTime)
	assert.Equal(t, nextKey.GotId.AasciValue, "a")
	assert.Equal(t, nextKey.NoteId.AasciValue, "rt") //this was incremented
	assert.Equal(t, engine.TimeToMillisString(nextKey.CreatedTime), nextMillis)

}
