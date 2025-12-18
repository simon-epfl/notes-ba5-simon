package rpc

type CreateRequest struct {
	Filename      string
	IsReplication bool
}

type CreateResponse struct {
	FileID     string
	Success    bool
	Error      string
	Generation int
}

type ReadRequest struct {
	FileID string
	Offset int64
	Length int
}

type ReadResponse struct {
	Data    []byte
	Success bool
	Error   string
}

type OpenRequest struct {
	Filename      string
	IsReplication bool
}

type OpenResponse struct {
	Success bool
	FileID  string
	Content []byte
	Error   string
}

type TestAuthRequest struct {
	FileID string
}

type TestAuthResponse struct {
	Valid bool
	Error string
}

type WriteRequest struct {
	FileID  string
	Content []byte
}

type WriteResponse struct {
	Success       bool
	Written       int
	NewGeneration int
	Error         string
}

type CloseRequest struct {
	FileID        string
	Content       []byte
	IsReplication bool
}

type CloseResponse struct {
	Success   bool
	Error     string
	NewFileID string
}

type PingRequest struct{}

type PingResponse struct {
	Alive     bool
	Timestamp int64
}

type GetFileListRequest struct {
	Directory string
}

type GetFileListResponse struct {
	Files   []string
	Success bool
	Error   string
}

type Operation struct {
	OpID       int64  `json:"op_id"`
	OpType     string `json:"op_type"`
	Filename   string `json:"filename"`
	Content    []byte `json:"content,omitempty"`
	IsOutput   bool   `json:"is_output"`
	Generation int    `json:"generation,omitempty"`
}

type ReplicateRequest struct {
	Op Operation
}

type ReplicateResponse struct {
	Success bool
	Error   string
}

type GetOperationsRequest struct {
	SinceOpID int64
}

type GetOperationsResponse struct {
	Operations []Operation
	Success    bool
	Error      string
}
