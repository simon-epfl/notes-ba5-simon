package rpc

import (
	assign "ds-uoe-vash/prime-app/assign"
	"time"
)

// all messagess used by worker and master for request and response
// currentyl I dont know what i need and what i dont need
// the struct names are self explanatory

// register request struct
type RegisterReq struct {
	WorkerID string
	Addr     string
}

// register request response
type RegisterRes struct {
	WorkerID     string
	UnixTime     int64
	SessionToken string
}

type HeartbeatReq struct {
	WorkerID string
	Load     int
}

type HeartbeatRes struct {
	UnixTime int64
	Commands []string
}

type RequestTaskReq struct {
	WorkerID string
}

type Assignment struct {
	TaskPath     string
	StartByte    int
	EndByte      int
	WorkerID     string
	AssignmentID string
	Assigned     bool
	Completed    bool
	Lease        time.Time

	// Lease        int64
	Attempt int
}

type RequestTaskRes struct {
	Assignment *Assignment
	NoTask     bool
}

type ReportMergeReq struct {
	WorkerID string
	Filepath string
}

type ReportMergeRes struct {
	Ack bool
}

// type StartMergeReq struct {
// 	Start bool
// }

// type StartMergeRes struct {
// 	Ack bool
// }

type ExtendLeaseReq struct {
	WorkerID     string
	AssignmentID string
	ExtraSeconds int64
}

type ExtendLeaseRes struct {
	Granted bool
	NewUnix int64
}

type ReportResultReq struct {
	WorkerID     string
	AssignmentId string
	ResultId     string
	// ResultPath   string
	Success bool
	Attempt int
	// Summary      string
}

type ReportResultRes struct {
	Ack bool
}

type GlobalSnapshot struct {
	CoordState    *CoordinatorState
	WorkerStates  map[string]*WorkerSnapshot
	ChannelStates map[string][]interface{}
	UnixTime      int64
}

type WorkerSnapshot struct {
	WorkerID     string
	Assignment   *Assignment
	FilesToMerge []string
	numOfMerges  int
}

type SnapshotMarkerReq struct {
}

type SnapshotMarkerRes struct {
	Ack bool
}

type CoordinatorState struct {
	Workers              map[string]*WorkerState
	BigFileToAssignments map[string][]*assign.Assignment
	Assignments          map[string]*assign.Assignment
	WorkerFilesToMerge   map[string][]string
	Unassigned           []string
	FinalMerge           []string
	MergeStarted         bool
}

type WorkerState struct {
	Id         string
	Addr       string
	Capacity   int
	Load       int
	LastActive time.Time
}
