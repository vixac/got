package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
)

//VX:TODO write some examples and ask claude

/*
LongForm keys are going to be in the following format: "<gotId>:<noteId>:<createMillis>"
The got Id component is the aasci value of the got id. The noteId is the aasci value of the bullet id. Note ids are
unique and incrementing under a given gotId, but they themselves are not unique. We also shove the createMillis into the key
just because we have nowhere else to put it.
*/

type LongFormKey struct {
	NoteId      bullet_stl.BulletId
	GotId       GotId
	CreatedTime time.Time
}

func (k *LongFormKey) Next(time time.Time) LongFormKey {
	return LongFormKey{
		GotId:       k.GotId,
		NoteId:      k.NoteId.Next(),
		CreatedTime: time,
	}
}

func FirstNoteId() bullet_stl.BulletId {
	id, _ := bullet_stl.NewBulletIdFromInt(1000)
	return *id
}

func (k *LongFormKey) ToString() string {
	createdStr := TimeToMillisString(k.CreatedTime)
	return k.GotId.AasciValue + ":" + k.NoteId.AasciValue + ":" + createdStr
}

func NewLongFormKeyFromString(input string) (*LongFormKey, error) {
	split := strings.Split(input, ":")
	if len(split) != 3 {
		fmt.Printf("VX: this is not a valid longform key: %s\n", input)
		return nil, errors.New("Invalid longform key")
	}
	gotInput := split[0]
	noteInput := split[1]
	createdMillisInput := split[2]
	gotId, err := NewGotId(gotInput)
	if err != nil {
		return nil, err
	}
	bulletId, err := bullet_stl.NewBulletIdFromString(noteInput)
	if err != nil {
		return nil, err
	}

	createdTime, err := EpochMillisStringToDate(createdMillisInput)
	if err != nil {
		return nil, err
	}
	longForm := LongFormKey{
		GotId:       *gotId,
		NoteId:      *bulletId,
		CreatedTime: *createdTime,
	}
	return &longForm, nil
}

func TimeToMillisString(time time.Time) string {
	millis := time.UnixMilli()
	return strconv.FormatInt(millis, 10)
}

func EpochMillisStringToDate(millisStr string) (*time.Time, error) {

	millis, err := strconv.ParseInt(millisStr, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(0, millis*int64(time.Millisecond))
	return &t, nil
}
