package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	netrpc "net/rpc"
	"os"
	"path/filepath"
	"sync"

	"ds-uoe-vash/dfs/pkg/rpc"
)

type FileServer struct {
	inputDir  string
	outputDir string
	metadata  map[string]*FileMetadata
	metaMutex sync.RWMutex

	serverAddrs []string

	operations []rpc.Operation
	lastOpID   int64
	opMutex    sync.RWMutex

	metadataFile   string
	operationsFile string
}

type FileMetadata struct {
	Generation int
}

type FileHandle struct {
	Path       string
	IsOutput   bool
	Generation int
}

func NewFileServer(inputDir, outputDir string, serverAddrs []string) *FileServer {
	fs := &FileServer{
		inputDir:       inputDir,
		outputDir:      outputDir,
		metadata:       make(map[string]*FileMetadata),
		serverAddrs:    serverAddrs,
		operations:     make([]rpc.Operation, 0, 100),
		lastOpID:       0,
		metadataFile:   filepath.Join(outputDir, ".metadata.json"),
		operationsFile: filepath.Join(outputDir, ".operations.json"),
	}

	fs.loadMetadata()
	fs.loadOperations()

	fs.recoverFromPeers()

	return fs
}

func (fs *FileServer) encodeHandle(fh *FileHandle) (string, error) {
	data, err := json.Marshal(fh)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (fs *FileServer) decodeHandle(encoded string) (*FileHandle, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid file handle format")
	}

	var fh FileHandle
	if err := json.Unmarshal(data, &fh); err != nil {
		return nil, fmt.Errorf("corrupted file handle")
	}

	return &fh, nil
}

func (fs *FileServer) getFilePath(fh *FileHandle) string {
	if fh.IsOutput {
		return filepath.Join(fs.outputDir, fh.Path)
	}
	return filepath.Join(fs.inputDir, fh.Path)
}

func (fs *FileServer) getOrCreateMetadata(path string) *FileMetadata {
	fs.metaMutex.Lock()
	defer fs.metaMutex.Unlock()

	if meta, exists := fs.metadata[path]; exists {
		return meta
	}

	meta := &FileMetadata{
		Generation: 0,
	}
	fs.metadata[path] = meta
	return meta
}

func (fs *FileServer) markMetadataAsModified(path string) *FileMetadata {
	fs.metaMutex.Lock()
	defer fs.metaMutex.Unlock()

	meta := fs.metadata[path]
	if meta == nil {
		meta = &FileMetadata{}
		fs.metadata[path] = meta
	}

	meta.Generation++
	return meta
}

func (fs *FileServer) getMetadata(path string) *FileMetadata {
	fs.metaMutex.RLock()
	defer fs.metaMutex.RUnlock()
	return fs.metadata[path]
}

func (fs *FileServer) addOperation(op rpc.Operation) {
	fs.opMutex.Lock()
	defer fs.opMutex.Unlock()

	fs.lastOpID++
	op.OpID = fs.lastOpID

	if len(fs.operations) >= 100 {
		fs.operations = fs.operations[1:]
	}
	fs.operations = append(fs.operations, op)

	fs.persistOperations()
}

func (fs *FileServer) broadcastOperation(op rpc.Operation, isFromClient bool) {
	if !isFromClient {

		return
	}

	for _, addr := range fs.serverAddrs {
		go func(address string) {
			client, err := netrpc.Dial("tcp", address)
			if err != nil {
				return
			}
			defer client.Close()

			var resp rpc.ReplicateResponse
			client.Call("FileServer.ReplicateOperation", &rpc.ReplicateRequest{Op: op}, &resp)
		}(addr)
	}
}

func (fs *FileServer) loadMetadata() {
	data, err := os.ReadFile(fs.metadataFile)
	if err != nil {
		return
	}

	fs.metaMutex.Lock()
	defer fs.metaMutex.Unlock()
	json.Unmarshal(data, &fs.metadata)
}

func (fs *FileServer) persistMetadata() {
	fs.metaMutex.RLock()
	data, err := json.Marshal(fs.metadata)
	fs.metaMutex.RUnlock()

	if err != nil {
		return
	}
	os.WriteFile(fs.metadataFile, data, 0644)
}

func (fs *FileServer) loadOperations() {
	data, err := os.ReadFile(fs.operationsFile)
	if err != nil {
		return
	}

	var ops []rpc.Operation
	if err := json.Unmarshal(data, &ops); err != nil {
		return
	}

	fs.opMutex.Lock()
	defer fs.opMutex.Unlock()
	fs.operations = ops
	if len(ops) > 0 {
		fs.lastOpID = ops[len(ops)-1].OpID
	}
}

func (fs *FileServer) persistOperations() {
	data, err := json.Marshal(fs.operations)
	if err != nil {
		return
	}
	os.WriteFile(fs.operationsFile, data, 0644)
}

func (fs *FileServer) recoverFromPeers() {
	fs.opMutex.RLock()
	myLastOpID := fs.lastOpID
	fs.opMutex.RUnlock()

	for _, addr := range fs.serverAddrs {
		client, err := netrpc.Dial("tcp", addr)
		if err != nil {
			continue
		}

		var resp rpc.GetOperationsResponse
		err = client.Call("FileServer.GetOperationsSince", &rpc.GetOperationsRequest{SinceOpID: myLastOpID}, &resp)
		client.Close()

		if err != nil {
			continue
		}

		for _, op := range resp.Operations {
			fs.replayOperation(op)
		}
	}
}

func (fs *FileServer) replayOperation(op rpc.Operation) {
	switch op.OpType {
	case "open":
		fs.metaMutex.Lock()
		if _, exists := fs.metadata[op.Filename]; !exists {
			fs.metadata[op.Filename] = &FileMetadata{Generation: op.Generation}
		}
		fs.metaMutex.Unlock()
	case "create":
		filePath := filepath.Join(fs.outputDir, op.Filename)
		os.Create(filePath)
		fs.metaMutex.Lock()
		fs.metadata[filePath] = &FileMetadata{Generation: op.Generation}
		fs.metaMutex.Unlock()
	case "close":
		if len(op.Content) > 0 {
			filePath := filepath.Join(fs.outputDir, op.Filename)
			os.WriteFile(filePath, op.Content, 0644)
			fs.metaMutex.Lock()
			if meta, exists := fs.metadata[filePath]; exists {
				meta.Generation++
			}
			fs.metaMutex.Unlock()
		}
	}

	fs.opMutex.Lock()
	if op.OpID > fs.lastOpID {
		fs.lastOpID = op.OpID
		if len(fs.operations) >= 100 {
			fs.operations = fs.operations[1:]
		}
		fs.operations = append(fs.operations, op)
	}
	fs.opMutex.Unlock()

	fs.persistMetadata()
	fs.persistOperations()
}

func (fs *FileServer) Open(req *rpc.OpenRequest, resp *rpc.OpenResponse) error {
	filePath := filepath.Join(fs.inputDir, req.Filename)

	_, err := os.Stat(filePath)
	if err != nil {
		resp.Success = false
		resp.Error = "file does not exist"
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		resp.Success = false
		resp.Error = fmt.Sprintf("failed to read file: %v", err)
		return nil
	}

	meta := fs.getOrCreateMetadata(filePath)

	op := rpc.Operation{
		OpType:     "open",
		Filename:   filePath,
		IsOutput:   false,
		Generation: meta.Generation,
	}
	fs.addOperation(op)
	fs.persistMetadata()
	fs.broadcastOperation(op, !req.IsReplication)

	fh := &FileHandle{
		Path:       req.Filename,
		IsOutput:   false,
		Generation: meta.Generation,
	}

	encoded, err := fs.encodeHandle(fh)
	if err != nil {
		resp.Success = false
		resp.Error = "failed to create file handle"
		return nil
	}

	resp.FileID = encoded
	resp.Content = content
	resp.Success = true

	return nil
}

func (fs *FileServer) Create(req *rpc.CreateRequest, resp *rpc.CreateResponse) error {
	filePath := filepath.Join(fs.outputDir, req.Filename)

	if _, err := os.Stat(filePath); err == nil {

		meta := fs.getOrCreateMetadata(filePath)

		fh := &FileHandle{
			Path:       req.Filename,
			IsOutput:   true,
			Generation: meta.Generation,
		}

		encoded, err := fs.encodeHandle(fh)
		if err != nil {
			resp.Success = false
			resp.Error = "failed to create file handle"
			return nil
		}

		resp.FileID = encoded
		resp.Generation = meta.Generation
		resp.Success = true
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		resp.Success = false
		resp.Error = "failed to create file"
		return nil
	}
	file.Close()

	meta := &FileMetadata{
		Generation: 1,
	}

	fs.metaMutex.Lock()
	fs.metadata[filePath] = meta
	fs.metaMutex.Unlock()

	op := rpc.Operation{
		OpType:     "create",
		Filename:   req.Filename,
		IsOutput:   true,
		Generation: 1,
	}
	fs.addOperation(op)
	fs.persistMetadata()
	fs.broadcastOperation(op, !req.IsReplication)

	fh := &FileHandle{
		Path:       req.Filename,
		IsOutput:   true,
		Generation: meta.Generation,
	}

	encoded, err := fs.encodeHandle(fh)
	if err != nil {
		resp.Success = false
		resp.Error = "failed to create file handle"
		return nil
	}

	resp.FileID = encoded
	resp.Generation = meta.Generation
	resp.Success = true
	return nil
}

func (fs *FileServer) TestAuth(req *rpc.TestAuthRequest, resp *rpc.TestAuthResponse) error {

	fh, err := fs.decodeHandle(req.FileID)
	if err != nil {
		resp.Valid = false
		resp.Error = "invalid file handle"
		return nil
	}

	filePath := fs.getFilePath(fh)

	meta := fs.getMetadata(filePath)
	if meta == nil {

		resp.Valid = false
		resp.Error = "file not found or never opened"
		return nil
	}

	if meta.Generation == fh.Generation {

		resp.Valid = true
	} else {

		resp.Valid = false
		resp.Error = "file has been modified"
	}

	return nil
}

func (fs *FileServer) Read(req *rpc.ReadRequest, resp *rpc.ReadResponse) error {
	resp.Success = false
	resp.Error = "Read operation is not supported; clients read from their cached copies"
	return nil
}

func (fs *FileServer) Write(req *rpc.WriteRequest, resp *rpc.WriteResponse) error {
	resp.Success = false
	resp.Error = "Write operation is not supported; writes are handled during Close"
	return nil
}

func (fs *FileServer) Close(req *rpc.CloseRequest, resp *rpc.CloseResponse) error {

	fh, err := fs.decodeHandle(req.FileID)
	if err != nil {
		resp.Success = false
		resp.Error = "invalid file handle"
		return nil
	}

	if len(req.Content) > 0 {
		filePath := fs.getFilePath(fh)

		tempPath := filePath + ".tmp"

		if err := os.WriteFile(tempPath, req.Content, 0644); err != nil {
			resp.Success = false
			resp.Error = fmt.Sprintf("failed to write file: %v", err)
			return nil
		}

		if err := os.Rename(tempPath, filePath); err != nil {
			os.Remove(tempPath)
			resp.Success = false
			resp.Error = fmt.Sprintf("failed to commit file: %v", err)
			return nil
		}

		newMeta := fs.markMetadataAsModified(filePath)

		op := rpc.Operation{
			OpType:     "close",
			Filename:   fh.Path,
			Content:    req.Content,
			IsOutput:   fh.IsOutput,
			Generation: newMeta.Generation,
		}
		fs.addOperation(op)
		fs.persistMetadata()
		fs.broadcastOperation(op, !req.IsReplication)

		fh := &FileHandle{
			Path:       filePath,
			IsOutput:   false,
			Generation: newMeta.Generation,
		}
		encoded, err := fs.encodeHandle(fh)
		if err != nil {
			resp.Success = false
			resp.Error = "failed to create file handle"
			return nil
		}

		resp.NewFileID = encoded
	}

	resp.Success = true
	return nil
}

func (fs *FileServer) ReplicateOperation(req *rpc.ReplicateRequest, resp *rpc.ReplicateResponse) error {
	fs.replayOperation(req.Op)
	resp.Success = true
	return nil
}

func (fs *FileServer) GetOperationsSince(req *rpc.GetOperationsRequest, resp *rpc.GetOperationsResponse) error {
	fs.opMutex.RLock()
	defer fs.opMutex.RUnlock()

	var ops []rpc.Operation
	for _, op := range fs.operations {
		if op.OpID > req.SinceOpID {
			ops = append(ops, op)
		}
	}

	resp.Operations = ops
	resp.Success = true
	return nil
}
