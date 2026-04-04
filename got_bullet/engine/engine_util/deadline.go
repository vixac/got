package engine_util

import (
	"errors"
	"time"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func Deadline(deadline *engine.DateTime, state engine.GotState, now time.Time) (string, console.Token, error) {

	var displayDeadline = ""
	var deadlineToken console.Token = console.TokenSecondary{}
	//VX:TODO get this date wrangling out. Its business logic	//if theres a deadline and either its a group or its an active job
	if deadline != nil && state == engine.Active {

		//this "n" is not strongly typed and I feel bad.
		//handle all the special cases
		if deadline.Special == "n" {
			return "---Now---", console.TokenNow{}, nil
		}

		//if its not special, its assumed to be a normal deadline

		deadlineDate, err := deadline.ToDate()
		if err != nil {
			return "", deadlineToken, err
		}
		if deadlineDate == nil {
			return "", deadlineToken, errors.New("Missing deadline date.")
		}

		deadStr, spaceTime := console.HumanizeDate(time.Time(*deadlineDate), now)
		displayDeadline = deadStr
		deadlineToken = console.ToToken(spaceTime)
		return displayDeadline, deadlineToken, nil
	}
	return "", deadlineToken, nil
}
