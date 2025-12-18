package client

import (
	"fmt"
	"io"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	dfsrpc "ds-uoe-vash/dfs/pkg/rpc"
)

const (
	maxRetries   = 3
	retryTimeout = 100 * time.Millisecond
)

type DFSClient struct {
	cacheDir    string
	rpcClient   *rpc.Client
	serverAddrs []string
	currentIdx  int

	openFiles map[string]*CachedFile
	fileMutex sync.RWMutex
}

type CachedFile struct {
	RemotePath string
	LocalPath  string
	FileID     string
	Modified   bool
	mutex      sync.Mutex
}

func NewDFSClient(serverAddrs []string, cacheDir string) (*DFSClient, error) {
	if len(serverAddrs) == 0 {
		return nil, fmt.Errorf("no server addresses provided")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %v", err)
	}

	var client *rpc.Client
	var connErr error

	for i, addr := range serverAddrs {
		client, connErr = rpc.Dial("tcp", addr)
		if connErr == nil {
			return &DFSClient{
				serverAddrs: serverAddrs,
				currentIdx:  i,
				cacheDir:    cacheDir,
				rpcClient:   client,
				openFiles:   make(map[string]*CachedFile),
			}, nil
		}
	}

	return nil, fmt.Errorf("failed to connect to any server: %v", connErr)
}

func (c *DFSClient) tryNextServer() error {
	startIdx := c.currentIdx

	for i := 0; i < len(c.serverAddrs); i++ {
		nextIdx := (startIdx + i + 1) % len(c.serverAddrs)
		addr := c.serverAddrs[nextIdx]

		client, err := rpc.Dial("tcp", addr)
		if err == nil {
			if c.rpcClient != nil {
				c.rpcClient.Close()
			}
			c.rpcClient = client
			c.currentIdx = nextIdx
			return nil
		}
	}

	return fmt.Errorf("all servers unavailable")
}

func (c *DFSClient) callWithRetry(method string, args interface{}, reply interface{}) error {

	fmt.Printf("Calling RPC method: %s on server: %s\n", method, c.serverAddrs[c.currentIdx])
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryTimeout)
		}

		err := c.rpcClient.Call(method, args, reply)
		if err == nil {
			return nil
		}

		lastErr = err

		if strings.Contains(err.Error(), "connection") ||
			strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "broken pipe") {
			if failoverErr := c.tryNextServer(); failoverErr == nil {
				fmt.Printf("Switched to next server: %s\n", c.serverAddrs[c.currentIdx])
				continue
			}
		}
	}

	return fmt.Errorf("RPC call failed after %d attempts: %v", maxRetries, lastErr)
}

func (c *DFSClient) Open(filename string) (*os.File, error) {
	c.fileMutex.Lock()
	defer c.fileMutex.Unlock()

	localPath := filepath.Join(c.cacheDir, filename)

	if cached, exists := c.openFiles[filename]; exists {

		return os.OpenFile(cached.LocalPath, os.O_RDWR, 0644)
	}

	var needFetch = true
	var cachedFileID = ""

	if _, err := os.Stat(localPath); err == nil {

		metaPath := localPath + ".meta"
		if metaData, err := os.ReadFile(metaPath); err == nil {
			cachedFileID = strings.TrimSpace(string(metaData))

			if cachedFileID != "" {

				req := &dfsrpc.TestAuthRequest{
					FileID: cachedFileID,
				}
				resp := &dfsrpc.TestAuthResponse{}

				if err := c.callWithRetry("FileServer.TestAuth", req, resp); err == nil && resp.Valid {

					needFetch = false
				} else {
					needFetch = true
				}

			}
		}
	}

	var fileID string

	if needFetch {

		req := &dfsrpc.OpenRequest{Filename: filename}
		resp := &dfsrpc.OpenResponse{}

		if err := c.callWithRetry("FileServer.Open", req, resp); err != nil {
			return nil, fmt.Errorf("RPC Open failed: %v", err)
		}

		if !resp.Success {
			return nil, fmt.Errorf("server error: %s", resp.Error)
		}

		if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache subdirectory: %v", err)
		}

		if err := os.WriteFile(localPath, resp.Content, 0644); err != nil {
			return nil, fmt.Errorf("failed to write cached file: %v", err)
		}

		metaPath := localPath + ".meta"
		os.WriteFile(metaPath, []byte(resp.FileID), 0644)

		fileID = resp.FileID
	} else {
		fileID = cachedFileID
	}

	c.openFiles[filename] = &CachedFile{
		RemotePath: filename,
		LocalPath:  localPath,
		FileID:     fileID,
		Modified:   false,
	}

	return os.OpenFile(localPath, os.O_RDWR, 0644)
}

func (c *DFSClient) Create(filename string) (*os.File, error) {
	c.fileMutex.Lock()
	defer c.fileMutex.Unlock()

	localPath := filepath.Join(c.cacheDir, filename)

	req := &dfsrpc.CreateRequest{Filename: filename}
	resp := &dfsrpc.CreateResponse{}

	if err := c.callWithRetry("FileServer.Create", req, resp); err != nil {
		return nil, fmt.Errorf("RPC Create failed: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("server error: %s", resp.Error)
	}

	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache subdirectory: %v", err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local cache: %v", err)
	}

	c.openFiles[filename] = &CachedFile{
		RemotePath: filename,
		LocalPath:  localPath,
		FileID:     resp.FileID,
		Modified:   true,
	}

	return file, nil
}

func (c *DFSClient) Read(file *os.File, buf []byte) (int, error) {

	return file.Read(buf)
}

func (c *DFSClient) Write(file *os.File, data []byte) (int, error) {

	c.fileMutex.RLock()
	filename := filepath.Base(file.Name())
	cached := c.openFiles[filename]
	c.fileMutex.RUnlock()

	if cached == nil {
		return 0, fmt.Errorf("file not opened through DFS client")
	}

	cached.mutex.Lock()
	cached.Modified = true
	cached.mutex.Unlock()

	return file.Write(data)
}

func (c *DFSClient) Close(file *os.File) error {
	filename := filepath.Base(file.Name())

	c.fileMutex.Lock()
	cached, exists := c.openFiles[filename]
	if !exists {
		c.fileMutex.Unlock()
		return file.Close()
	}
	delete(c.openFiles, filename)
	c.fileMutex.Unlock()

	if err := file.Close(); err != nil {
		return err
	}

	var content []byte
	if cached.Modified {

		var err error
		content, err = os.ReadFile(cached.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to read local file for flush: %v", err)
		}
	}

	closeReq := &dfsrpc.CloseRequest{
		FileID:  cached.FileID,
		Content: content,
	}
	closeResp := &dfsrpc.CloseResponse{}

	if err := c.callWithRetry("FileServer.Close", closeReq, closeResp); err != nil {
		return fmt.Errorf("RPC Close failed: %v", err)
	}

	if !closeResp.Success {
		return fmt.Errorf("server close error: %s", closeResp.Error)
	}

	metaPath := cached.LocalPath + ".meta"
	if closeResp.NewFileID != "" {
		os.WriteFile(metaPath, []byte(closeResp.NewFileID), 0644)
	}

	return nil
}

func (c *DFSClient) ReadAt(file *os.File, buf []byte, offset int64) (int, error) {
	return file.ReadAt(buf, offset)
}

func (c *DFSClient) WriteAt(file *os.File, data []byte, offset int64) (int, error) {
	filename := filepath.Base(file.Name())

	c.fileMutex.RLock()
	cached := c.openFiles[filename]
	c.fileMutex.RUnlock()

	if cached == nil {
		return 0, fmt.Errorf("file not opened through DFS client")
	}

	cached.mutex.Lock()
	cached.Modified = true
	cached.mutex.Unlock()

	return file.WriteAt(data, offset)
}

func (c *DFSClient) Seek(file *os.File, offset int64, whence int) (int64, error) {
	return file.Seek(offset, whence)
}

func (c *DFSClient) Stat(file *os.File) (os.FileInfo, error) {
	return file.Stat()
}

func (c *DFSClient) Disconnect() error {
	if c.rpcClient != nil {
		return c.rpcClient.Close()
	}
	return nil
}

func (c *DFSClient) ReadFull(file *os.File) ([]byte, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size())
	_, err = io.ReadFull(file, buf)
	return buf, err
}

func (c *DFSClient) WriteAll(file *os.File, data []byte) error {

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	_, err := c.Write(file, data)
	return err
}
