package coordinator

import (
	"bufio"
	"bytes"
	"container/heap"
	"os"
	"runtime"
	"strconv"

	// "bytes"
	dfsClient "ds-uoe-vash/dfs/pkg/client"
	minheap "ds-uoe-vash/prime-app/heapStruct"
	"ds-uoe-vash/prime-app/rpc"
	primerpc "ds-uoe-vash/prime-app/rpc"

	"fmt"
	legitRpc "net/rpc"
	"strings"

	// "io"
	"log"
	// "os"
	assign "ds-uoe-vash/prime-app/assign"
	"encoding/json"
	fp "path/filepath"
	"sync"
	"time"
)

// maybe some of the maps are temporary, not quite sure what we need yet

type Coordinator struct {
	mutex   sync.Mutex //prevents clients from accessing at the same time to preven data races
	dfsRoot string     // root fo dfs file system where we start working on

	Workers map[string]*primerpc.WorkerState // current active workers

	workerToAssign map[string]*assign.Assignment // current activ assignemnts

	bigFileToAssignments map[string][]*assign.Assignment //Tracking if all the small chunks of the file have been completed

	assignments map[string]*assign.Assignment // map from assignment path to assignmentobject for o(1) lookup

	unassigned []string //List of available assignments in a FIFO manner

	FinalMerge []string

	workerFilesToMerge map[string][]string

	Recovered bool

	timeout int

	MergeStarted bool

	dfsClient              *dfsClient.DFSClient
	SnapshotInProg         bool
	snapshotMarkerReceived map[string]bool
	workerSnapshots        map[string]*primerpc.WorkerSnapshot
	messagesInTransit      map[string][]interface{}
}

func ConstructCoordinator(dfsIP string, dfsRoot string, timeout int) *Coordinator {
	client, err := dfsClient.NewDFSClient([]string{dfsIP}, "./cache") // chnged dfsIP
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect()

	coordinator := &Coordinator{
		dfsRoot:              dfsRoot,
		Workers:              make(map[string]*primerpc.WorkerState), // string = workerID
		bigFileToAssignments: make(map[string][]*assign.Assignment),  //Haaroon check out (map[string][]*Assignment)// string = filename
		workerToAssign:       make(map[string]*assign.Assignment),    // string = workerID
		assignments:          make(map[string]*assign.Assignment),
		workerFilesToMerge:   make(map[string][]string), // string = filename
		timeout:              timeout,
		Recovered:            false,
		dfsClient:            client,
	}

	err = coordinator.RecoverFromSnapshot()
	if err != nil {
		log.Printf("Could not recover from existing snapshot")
	} else {
		log.Printf("Coord recovered from snapshot")
	}

	// go coordinator.loop()
	return coordinator
}

// Register a worker server.
func (c *Coordinator) RegisterWorker(req *primerpc.RegisterReq, resp *primerpc.RegisterRes) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.SnapshotInProg && c.snapshotMarkerReceived[req.WorkerID] == false {
		c.messagesInTransit[req.WorkerID] = append(c.messagesInTransit[req.WorkerID], *req)
	}
	log.Printf("COORD says Registered worker %s", req.WorkerID)
	c.Workers[req.WorkerID] = &primerpc.WorkerState{
		Id:         req.WorkerID,
		Addr:       req.Addr,
		LastActive: time.Now(),
	}

	resp.WorkerID = req.WorkerID
	resp.UnixTime = time.Now().Unix()
	return nil
}

func (c *Coordinator) ReportMerge(req *primerpc.ReportMergeReq, resp *primerpc.ReportMergeRes) error {
	c.mutex.Lock()
	if c.SnapshotInProg && c.snapshotMarkerReceived[req.WorkerID] == false {
		c.messagesInTransit[req.WorkerID] = append(c.messagesInTransit[req.WorkerID], *req)
	}

	log.Printf("COORD is getting reported a merge")
	if req.Filepath != "" {
		c.FinalMerge = append(c.FinalMerge, req.Filepath)
		log.Printf("COORD added file %s to finalmerge list", req.Filepath)

	}
	resp.Ack = true
	done := len(c.FinalMerge) == len(c.workerFilesToMerge)
	defer c.mutex.Unlock()

	if done && c.MergeStarted == false {

		go func() {
			filesGood := c.checkMergeFiles(30 * time.Second)

			if !filesGood {
				log.Printf("Only partial full output, files not ready")
			}

			log.Println("Files good, starting final merge")
			err := c.FinalMergeFunc(64*1024, 1024*1024)
			if err != nil {
				log.Printf("Final merge failed somewhere figure it out")
			} else {
				log.Printf("Lesgooo, final merge is DONE")
			}

		}()
	}
	return nil
}

func (c *Coordinator) checkMergeFiles(timeout time.Duration) bool {

	start := time.Now()
	//do busy wait to  wait for full file from server, this was hard to figure out
	for {
		allReady := true
		for _, filepath := range c.FinalMerge {
			file, err := c.dfsClient.Open(filepath)
			if err != nil {
				log.Printf("file %s cannot be opened %v ", filepath, err)
				allReady = false
				continue
			}

			scanner := bufio.NewScanner(file)
			nonEmpty := scanner.Scan()
			c.dfsClient.Close(file)

			if !nonEmpty {
				log.Printf("file %s is still empty", filepath)
				rel := fp.Base(filepath)
				cachePath := fp.Join("./cache", rel)
				metapath := cachePath + ".meta"
				os.Remove(cachePath)
				os.Remove(metapath)
				allReady = false
			}
		}
		if allReady {
			log.Printf("All files ready to be merged")
			return true
		}

		if time.Since(start) > timeout {
			log.Printf("Took too long to get files TIMEOUT")
			return false
		}
		time.Sleep(2 * time.Second)
	}
}

// func (c *Coordinator) checkMergeFiles() bool {
// 	for _, filepath := range c.FinalMerge {
// 		file, err := c.dfsClient.Open(filepath)
// 		if err != nil {
// 			log.Printf("File not able to be opened %s, %v", filepath, err)
// 			return false
// 		}
// 		scanner := bufio.NewScanner(file)
// 		nonEmpty := scanner.Scan()
// 		c.dfsClient.Close(file)

// 		if !nonEmpty {
// 			log.Printf("file %s is empty", filepath)
// 			return false
// 		}
// 	}
// 	return true
// }

// probably useless since we have snapshot
/*func (c *Coordinator) Heartbeat(req *primerpc.HeartbeatReq, resp *primerpc.HeartbeatRes) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	worker, exists := c.Workers[req.WorkerID]
	if exists {
		worker.LastActive = time.Now()
		worker.Load = req.Load
	} else {
		log.Printf("COORd: Received a heardbeat from %s", req.WorkerID)
	}
	resp.UnixTime = time.Now().Unix()
	return nil
}*/

func (c *Coordinator) RequestTask(req *primerpc.RequestTaskReq, resp *primerpc.RequestTaskRes) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.SnapshotInProg && c.snapshotMarkerReceived[req.WorkerID] == false {
		c.messagesInTransit[req.WorkerID] = append(c.messagesInTransit[req.WorkerID], *req)
	}

	task := c.FindTask(req.WorkerID)
	if task != nil {
		rpcAssignment := &rpc.Assignment{
			TaskPath:     task.TaskPath,
			WorkerID:     req.WorkerID,
			AssignmentID: task.AssignmentID,
			StartByte:    task.StartByte,
			EndByte:      task.EndByte,
			Lease:        task.Lease,
			Assigned:     true,
			Completed:    false,
		}

		resp.Assignment = rpcAssignment
		log.Printf("COORD Assigned task %s to worker %s", task.AssignmentID, task.WorkerID)
	} else {
		resp.Assignment = nil
		//log.Printf("COORD No tasks available")
	}
	return nil
}

func (c *Coordinator) ReportResult(req *primerpc.ReportResultReq, resp *primerpc.ReportResultRes) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.SnapshotInProg && c.snapshotMarkerReceived[req.WorkerID] == false {
		c.messagesInTransit[req.WorkerID] = append(c.messagesInTransit[req.WorkerID], *req)
	}

	task, exists := c.assignments[req.AssignmentId]
	if !exists {
		log.Printf("ERROR: Received report for unknown task: %s", req.AssignmentId)
		return fmt.Errorf("unknown assignment ID")
	}

	if req.Success {
		task.Completed = true
		c.workerFilesToMerge[req.WorkerID] = append(c.workerFilesToMerge[req.WorkerID], req.ResultId)
		log.Printf("Task completed successfully: %s by worker %s", req.AssignmentId, req.WorkerID)
	} else {
		task.Assigned = false
		task.WorkerID = ""
		c.unassigned = append(c.unassigned, task.AssignmentID)
		log.Printf("Task failed: %s. Re-queueing.", req.AssignmentId)
	}

	delete(c.workerToAssign, req.WorkerID)
	resp.Ack = true
	return nil
}

/*func (c *Coordinator) MonitorWorkers() {
	for {
		time.Sleep(30 * time.Second)

		c.mutex.Lock()

		timeout := 60 * time.Second
		var deadWorkers []string

		for workerID, workerState := range c.Workers {
			if time.Since(workerState.LastActive) > timeout {
				deadWorkers = append(deadWorkers, workerID)
			}
		}

		for _, workerID := range deadWorkers {
			if task, ok := c.workerToAssign[workerID]; ok {
				task.Assigned = false
				task.WorkerID = ""
				c.unassigned = append(c.unassigned, task.TaskPath)
				delete(c.workerToAssign, workerID)
			}
			delete(c.Workers, workerID)
		}
		c.mutex.Unlock()
	}
}*/

// Helper funtion used in RequestTask. Finds a task from the unassigned and
// writes its workerID in the assignment struct

func (c *Coordinator) FindTask(workerID string) *assign.Assignment {

	if len(c.unassigned) == 0 {
		return nil
	}
	taskPath := c.unassigned[0]
	c.unassigned = c.unassigned[1:]

	task := c.assignments[taskPath]
	task.WorkerID = workerID
	task.Assigned = true
	task.Lease = time.Now()
	return task
}

// linesPerFile is currently useless, it is stuck from last version
// chunkSize is the size of the chunk in bytes. 50000000 = 50MB file
func (c *Coordinator) ChunkFile(filepath string, linesPerFile int, chunkSize int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//optional later to fix the buffer size of the scanner, not the default 64KB

	file, err := c.dfsClient.Open(filepath) //open big file
	if err != nil {
		return fmt.Errorf("failed to open file in DFS: %v", err)
	}
	defer c.dfsClient.Close(file)

	scanner := bufio.NewScanner(file) //create scanner to read the lines

	currentOff := 0
	chunkStartOff := 0
	chunkNr := 0

	c.bigFileToAssignments[filepath] = []*assign.Assignment{} //initialized to be used in the map

	outName := strings.TrimSuffix(filepath, fp.Ext(filepath))
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		lineLength := len(lineBytes) + 1 //also add a byte for the \n

		if currentOff-chunkStartOff >= chunkSize {
			newFunction1(outName, chunkNr, filepath, chunkStartOff, currentOff, c)

			chunkNr++
			chunkStartOff = currentOff
		}
		currentOff += lineLength

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("we got a scan error boysss: %v", err)
	}

	if chunkStartOff < currentOff { //handle the last file
		newFunction1(outName, chunkNr, filepath, chunkStartOff, currentOff, c)

	}
	log.Printf(("created all chunks"))
	runtime.GC()
	return nil
}

func newFunction1(outName string, chunkNr int, filepath string, chunkStartOff int, currentOff int, c *Coordinator) {
	assignmentID := fmt.Sprintf("%s_%d", outName, chunkNr)
	newAssignment := &assign.Assignment{
		TaskPath:     filepath,
		StartByte:    chunkStartOff,
		EndByte:      currentOff,
		AssignmentID: assignmentID,
		Assigned:     false,
		Completed:    false,
		Lease:        time.Now(),
	}
	log.Printf("COORD: Created chunk %d: bytes %d to %d", chunkNr, chunkStartOff, currentOff)
	c.bigFileToAssignments[filepath] = append(c.bigFileToAssignments[filepath], newAssignment)
	c.assignments[assignmentID] = newAssignment
	c.unassigned = append(c.unassigned, assignmentID)
}

// refactored a helper fn to create assignment and add it to maps and lists

// Same function as in worker but modified to fit coordinator for final merge
type reader struct {
	f       *os.File
	s       *bufio.Scanner
	current uint64
}

func (c *Coordinator) FinalMergeFunc(bufferSize int, maxScanSize int) error {

	c.MergeStarted = true

	if len(c.FinalMerge) == 0 {
		return fmt.Errorf("Worker has no files to merge")
	}

	out := "primes.txt"
	outFile, err := c.dfsClient.Create(out)
	if err != nil {
		log.Printf("MergeFile creation failed: %v", err)
	}

	defer c.dfsClient.Close(outFile)

	writer := bufio.NewWriterSize(outFile, bufferSize)

	readers := make([]*reader, 0, len(c.FinalMerge))

	for _, filepath := range c.FinalMerge {
		file, err := c.dfsClient.Open(filepath)
		if err != nil {

			return fmt.Errorf("Could not open in worker merger: %v", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, bufferSize), maxScanSize)
		reader := &reader{
			f: file,
			s: scanner,
		}
		readers = append(readers, reader)
	}

	h := &minheap.MinHeap{}
	heap.Init(h)
	for i, r := range readers {
		if r.s.Scan() {
			val, err := strconv.ParseUint(r.s.Text(), 10, 64)
			if err != nil {
				return fmt.Errorf("uint parsing failed for worker in file %v whit err %v", i, err)
			}

			r.current = val
			toPush := minheap.KthInt{
				Val: val,
				Idx: i,
			}
			heap.Push(h, toPush)
		}
	}

	var last uint64
	first := true

	for h.Len() > 0 {
		item := heap.Pop(h).(minheap.KthInt)
		r := readers[item.Idx]
		v := item.Val

		if first || v != last {
			str := strconv.FormatUint(v, 10)
			_, _ = writer.WriteString(str) //avoid writng duplicates
			writer.WriteByte('\n')
			last = v
			first = false
		}

		if r.s.Scan() {
			line := r.s.Text()
			nextVal, err := strconv.ParseUint(line, 10, 64)

			if err != nil {
				return fmt.Errorf("uint parsing failed for coordinator in file %d whit err %v", item.Idx, err)
			}

			r.current = nextVal
			a := minheap.KthInt{
				Val: nextVal,
				Idx: item.Idx,
			}
			heap.Push(h, a)

		}
	}

	err1 := writer.Flush()
	if err1 != nil {
		return fmt.Errorf("Flushing at the end failes")
	}

	for _, r := range readers {
		c.dfsClient.Close(r.f)
	}
	log.Printf("Coordiantor has done the final merge")

	c.FinalMerge = c.FinalMerge[:0]
	return nil
}

func (c *Coordinator) ChunkFileSeparateFiles(filepath string, linesPerFile int, chunkSize int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//optional later to fix the buffer size of the scanner, not the default 64KB

	baseFilename := fp.Base(filepath)

	file, err := c.dfsClient.Open(filepath) //open big file
	if err != nil {
		return fmt.Errorf("failed to open file in DFS: %v", err)
	}
	defer c.dfsClient.Close(file) //defer closing until end

	partCount := 1
	lineCount := 0

	var outFile *os.File

	var outName string
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file) //create scanner to read the lines

	c.bigFileToAssignments[filepath] = []*assign.Assignment{} //initialized to be used in the map
	outNameBase := strings.TrimSuffix(baseFilename, fp.Ext(baseFilename))
	outName = fmt.Sprintf("%s_part_%d.txt", outNameBase, partCount)
	outFile, err = c.dfsClient.Create(outName) // create a file
	if err != nil {
		return fmt.Errorf("create failed brooo: %v", err)
	}

	for scanner.Scan() {
		buf.Write((scanner.Bytes()))
		buf.WriteByte('\n') //need to add the newline again because scanner removes it
		lineCount++

		if lineCount >= linesPerFile {
			err := c.dfsClient.WriteAll(outFile, buf.Bytes())
			if err != nil { // write it to a file
				return fmt.Errorf("Could not wriite file %v", err)
			}
			c.dfsClient.Close(outFile)        //close the file after we wrote in it
			newFunction(outName, c, filepath) //create the assignment object and fill in the studd

			partCount++   //increase part count
			lineCount = 0 //reset both line count and buffer to get ready for next file
			buf.Reset()
			//start the scanner again to check if we really need a new file. Start the loop a bit early but we still keep it in the buffer.
			if scanner.Scan() { //there was a problem that it created an extra empty file
				buf.Write((scanner.Bytes())) //this is a sketchy quick fix, could probably fix it later in a better way
				buf.WriteByte('\n')
				lineCount++
				outName = fmt.Sprintf("%s_part_%d.txt", outNameBase, partCount) //create the next file
				outFile, err = c.dfsClient.Create(outName)
				if err != nil {
					return fmt.Errorf("create failed: %v", err)
				}
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("we got a scan error boysss: %v", err)
	}

	if buf.Len() > 0 { //handle the last file
		if err := c.dfsClient.WriteAll(outFile, buf.Bytes()); err != nil {
			return fmt.Errorf("Final write failed : %v", err)
		}
		c.dfsClient.Close(outFile)
		newFunction(outName, c, filepath)
	} else {
		c.dfsClient.Close(outFile) //close if no new lines
	}

	log.Printf(("created all files"))
	return nil
}

// refactored a helper fn to create assignment and add it to maps and lists
func newFunction(outName string, c *Coordinator, filepath string) {
	newAssignment := &assign.Assignment{
		TaskPath:     outName,
		WorkerID:     "",
		AssignmentID: outName,
		Assigned:     false,
		Completed:    false,
		Lease:        time.Now(),
	}
	c.bigFileToAssignments[filepath] = append(c.bigFileToAssignments[filepath], newAssignment)
	c.assignments[outName] = newAssignment
	c.unassigned = append(c.unassigned, outName)
}

func (c *Coordinator) ChunkFiles(filepath string, numlines int, chunksize int) error {
	if numlines < 1 {
		return c.ChunkFile(filepath, numlines, chunksize)
	} else {
		return c.ChunkFileSeparateFiles(filepath, numlines, chunksize)
	}
}

func (c *Coordinator) getLocal() *primerpc.CoordinatorState {
	state := &primerpc.CoordinatorState{
		Workers:              c.Workers,
		BigFileToAssignments: c.bigFileToAssignments,
		Assignments:          c.assignments,
		Unassigned:           c.unassigned,
		FinalMerge:           c.FinalMerge,
		MergeStarted:         c.MergeStarted,
		WorkerFilesToMerge:   c.workerFilesToMerge,
	}
	return state
}

func (c *Coordinator) StartSnapshot() error {
	c.mutex.Lock()
	if c.SnapshotInProg {
		c.mutex.Unlock()
		return fmt.Errorf("snapshot already in progress")
	}
	c.SnapshotInProg = true
	c.snapshotMarkerReceived = make(map[string]bool)
	c.messagesInTransit = make(map[string][]interface{})
	c.workerSnapshots = make(map[string]*primerpc.WorkerSnapshot)

	_ = c.getLocal()
	for workerID, workerState := range c.Workers {
		c.messagesInTransit[workerID] = make([]interface{}, 0)
		go func(id string, addr string) {
			client, err := legitRpc.Dial("tcp", addr)
			if err != nil {
				log.Printf("Failed to dial worker")
				return
			} else {
				log.Printf("Dialed %s, ", addr)
			}
			defer client.Close()

			req := &primerpc.SnapshotMarkerReq{}
			resp := &primerpc.SnapshotMarkerRes{}
			err = client.Call("Worker.ReceiveMarker", req, resp)
			if err != nil {
				log.Printf("Failed to send marker RPC to worker %s: %v", id, err)
			} else {
				log.Printf("Sent marker, %s, %s", id, addr)
			}
		}(workerID, workerState.Addr)
	}
	c.mutex.Unlock()
	return nil
}

func (c *Coordinator) CollectWorkerSnapshot(req *primerpc.WorkerSnapshot, resp *primerpc.ReportResultRes) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	workerID := req.WorkerID
	if !c.SnapshotInProg {
		resp.Ack = true
		return fmt.Errorf("no Snapshot")
	}

	c.snapshotMarkerReceived[workerID] = true
	c.workerSnapshots[workerID] = req

	if len(c.workerSnapshots) == len(c.Workers) {
		log.Println("Global Snapshot Complete!")
		c.persistGlobalSnapshot()
		c.SnapshotInProg = false
	}
	resp.Ack = true
	return nil
}

func (c *Coordinator) persistGlobalSnapshot() {
	snapshot := &primerpc.GlobalSnapshot{
		CoordState:    c.getLocal(),
		WorkerStates:  c.workerSnapshots,
		ChannelStates: c.messagesInTransit,
		UnixTime:      time.Now().Unix(),
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		log.Printf("Could nto write snapshot to Json %v", err)
		return
	}

	filename := "snapshot.json"

	outFile, err := c.dfsClient.Create(filename)
	if err != nil {
		log.Printf("creating snapshot file failed %s in DFS: %v", filename, err)
		return
	}
	defer c.dfsClient.Close(outFile)

	if err := c.dfsClient.WriteAll(outFile, data); err != nil {
		log.Printf("Nahh failed to write snapshot DFS: %v", err)
		return
	}

	log.Printf("Global snapshot sent to dfs aa %s", filename)
	for id, worker := range c.Workers {
		go func(workerID string, addr string) {
			client, err := legitRpc.Dial("tcp", addr)
			if err != nil {
				log.Printf("Failed to dial worker %s for reset: %v", workerID, err)
				return
			}
			defer client.Close()

			req := &primerpc.SnapshotMarkerReq{}
			resp := &primerpc.SnapshotMarkerRes{}
			err = client.Call("Worker.ResetSnapshotState", req, resp)
			if err != nil {
				log.Printf("Failed to reset worker %s: %v", workerID, err)
			}
		}(id, worker.Addr)
	}
}

func (c *Coordinator) RecoverFromSnapshot() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := c.dfsClient.Open("snapshot.json")
	if err != nil {
		log.Printf("There is no snapshot")
		return nil
	}

	defer c.dfsClient.Close(file)

	var snapshot primerpc.GlobalSnapshot

	err = json.NewDecoder(file).Decode(&snapshot)
	if err != nil {
		return fmt.Errorf("COORD Could not decode that json snapshot maaan ")
	}

	if snapshot.CoordState == nil {
		return fmt.Errorf("Somehow snapshot has no coord stuff")
	}

	c.Workers = snapshot.CoordState.Workers
	c.bigFileToAssignments = snapshot.CoordState.BigFileToAssignments
	c.assignments = snapshot.CoordState.Assignments
	c.unassigned = snapshot.CoordState.Unassigned
	c.FinalMerge = snapshot.CoordState.FinalMerge
	c.MergeStarted = snapshot.CoordState.MergeStarted
	c.workerFilesToMerge = snapshot.CoordState.WorkerFilesToMerge
	c.Recovered = true

	c.workerToAssign = make(map[string]*assign.Assignment)
	for _, assigned := range c.assignments {
		if assigned.WorkerID != "" && !assigned.Completed {
			c.workerToAssign[assigned.WorkerID] = assigned
		}
	}

	c.SnapshotInProg = false
	c.snapshotMarkerReceived = make(map[string]bool)
	c.workerSnapshots = make(map[string]*primerpc.WorkerSnapshot)
	c.messagesInTransit = make(map[string][]interface{})

	log.Printf("Coord recovered from Snapshoooot")
	return nil
}
