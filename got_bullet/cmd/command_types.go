package cmd

import "time"

type Gid struct {
	Id string
}
type Deadline struct {
	Absolute         time.Time
	DaysBeforeParent int
}
type AddCmd struct {
	Under    *Gid
	Deadline *Deadline
}

type JobsCmd struct {
	Under *Gid
}
