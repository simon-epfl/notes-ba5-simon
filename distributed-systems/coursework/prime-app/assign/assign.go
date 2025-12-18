package coordinator

import "time"

// I am hoping to implement some worker state for load management maybe, working on that, also may scrap it.

// currently I have two assignment structs,
// it wasnt working before for some reason becuase i think they were in the same package so i needed a new one
// will remove it later so we only have one
type Assignment struct {
	TaskPath     string
	StartByte    int
	EndByte      int
	WorkerID     string
	AssignmentID string
	Assigned     bool
	Completed    bool
	Lease        time.Time
}
