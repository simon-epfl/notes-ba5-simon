package worker

import (
	// "fmt"
	"bufio"
	"container/heap"
	"encoding/json"
	"runtime"
	"slices"
	"strings"
	"time"

	// "strings"

	// "bytes"
	dfsClient "ds-uoe-vash/dfs/pkg/client"
	minheap "ds-uoe-vash/prime-app/heapStruct"
	primerpc "ds-uoe-vash/prime-app/rpc"
	fp "path/filepath"

	"fmt"
	"log"
	"math/big"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

type Worker struct {
	WorkerID      string
	address       string
	coordAddr     string
	Assignment    *primerpc.Assignment
	dfsClient     *dfsClient.DFSClient
	coordClient   *rpc.Client
	filesToMerge  []string
	numOfMerges   int
	stateRecorded bool
	mutex         sync.Mutex
}

// Worker constructor, also registers the worker to a coordinator.
// Cannot create constructor without a coordinator
func ConstructWorker(id string, dfsIP string, addr string, coordAddr string) *Worker {

	client, err := dfsClient.NewDFSClient([]string{dfsIP}, "./cache")
	if err != nil {
		log.Fatal(err)
	}
	coordClient, err := rpc.Dial("tcp", coordAddr)
	if err != nil {
		log.Fatal(err)
	}

	worker := &Worker{
		WorkerID:    id,
		address:     addr,
		coordAddr:   coordAddr,
		dfsClient:   client,
		coordClient: coordClient,
		numOfMerges: 0,
	}
	//register the worker with a cooridinator, would not make sense for a worker to exist without coordinator

	err = worker.snapshotRecovery()
	if err != nil {
		log.Printf("Could not recove worker %s, %v", worker.WorkerID, err)
	} else {
		log.Printf("recovered %s", worker.WorkerID)
	}

	err1 := worker.register()
	if err1 != nil {
		log.Printf("sadly could not register worker %s : %v", worker.WorkerID, err1)
		return nil
	}
	return worker
}

//I have not implemented the prime yet, working on that, also need the working loop and heartbeat loop.

// Register the worker to the coordinator
func (w *Worker) register() error {
	req := &primerpc.RegisterReq{
		WorkerID: w.WorkerID,
		Addr:     w.address,
	}
	resp := &primerpc.RegisterRes{}
	err := w.coordClient.Call("Coordinator.RegisterWorker", req, resp)
	if err != nil {

		return err
	}

	log.Printf("Registered worker %s.", w.WorkerID)
	return nil
}

// request and receive task from coordinator.
// calls the RPC function from Coord.
func (w *Worker) RequestTask() error {
	req := &primerpc.RequestTaskReq{
		WorkerID: w.WorkerID,
	}
	resp := &primerpc.RequestTaskRes{}

	err := w.coordClient.Call("Coordinator.RequestTask", req, resp)
	if err != nil {
		if connShutDown(err) {
			log.Printf("conn closed trying again")
			recErr := w.reconnect()
			if recErr != nil {
				log.Printf("Could never reconnect worker %s,", w.WorkerID)
			} else {
				log.Printf("Reconnected wroker %s,", w.WorkerID)
			}
		} else {
			log.Printf("sadly could not get task for worker %s : %v", w.WorkerID, err)
		}
		return nil
	}
	if resp.Assignment == nil {
		//log.Printf("No tasks availbale")
		return nil
	}
	log.Printf("Worker %s Received task %s", w.WorkerID, resp.Assignment.AssignmentID)
	w.Assignment = resp.Assignment
	return nil
}

// func (w *Worker) reportResult2() error {
// 	req := &primerpc.ReportResultReq{
// 		WorkerID:     w.workerID,
// 		AssignmentId: w.Assignment.AssignmentID,
// 		Success:      true,
// 		Attempt:      w.Assignment.Attempt,
// 	}
// 	resp := &primerpc.ReportResultRes{}
// 	err := w.coordClient.Call("Coordinator.ReportTask", req, resp)
// 	log.Printf("Worker %s reported assignment %s", w.workerID, w.Assignment.AssignmentID)
// 	if err != nil {
// 		log.Printf("Whoops, Worker %s could not report result for assignment %s", w.workerID, w.Assignment.AssignmentID)
// 		return nil
// 	}
// 	if resp.Ack {
// 		log.Printf("Worker %s Got ack for %s ", w.workerID, w.Assignment.AssignmentID)
// 	}
// 	return nil
// }

// Check if a number is prime using Go's ProbablyPrime function Miller-Rabin.
// It is deterministic for 64 bit int.
func isPrime(x uint64) bool {
	n := new(big.Int).SetUint64(x)
	return n.ProbablyPrime(0)
}

// The Function that finds all primes with no duplicates then writes them to a new file on server.
func (w *Worker) PrimeWriteToFile() (string, error) {

	assignment := w.Assignment //get assignment from the worker object
	if assignment == nil {
		return "", fmt.Errorf("No assignment available ")
	}
	filepath := assignment.TaskPath //get the filename / path

	file, err := w.dfsClient.Open(filepath)
	if err != nil {

		return "", fmt.Errorf("failed to open file to check for primes: %v", err)
	} //open the assignment file
	defer w.dfsClient.Close(file)

	localMap := make(map[uint64]bool) //create a map to ensure unique primes
	//reduces load on coordinator in the merging step

	var outFile *os.File

	outName := assignment.AssignmentID
	outName = fmt.Sprintf("%s_out.txt", outName)

	if w.Assignment.StartByte > 0 {
		_, err := w.dfsClient.Seek(file, int64(assignment.StartByte), 0)
		if err != nil {
			return "", fmt.Errorf("Seek didnt work at %d: %v", assignment.StartByte, err)
		}
		log.Printf("Worker %s: seek worked", w.WorkerID)
	}

	outFile, err = w.dfsClient.Create(outName) //create new output file
	if err != nil {
		return "", fmt.Errorf("Could not create a file to write prime output: %v", err)
	}
	defer w.dfsClient.Close(outFile)

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)
	primecount := 0 //used to track nr of primes for debugging / logs
	currentByte := assignment.StartByte
	endByte := assignment.EndByte
	totalLines := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineLength := len(line) + 1

		if currentByte >= endByte {
			break
		}
		totalLines++

		nr, err := strconv.ParseUint(line, 10, 64) //need to convert line string to uint for prime func
		if err != nil {
			log.Printf("Failed to parse line: %q", line)
		}

		if isPrime(nr) {
			primecount++
			localMap[nr] = true

		}
		currentByte += lineLength
	}

	log.Printf("Worker %s Fully processed %d lines", w.WorkerID, totalLines)

	keys := make([]uint64, 0, len(localMap))
	for k := range localMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, nr := range keys {
		str := strconv.FormatUint(nr, 10) //need to convert to string for writer to write it
		writer.WriteString(str)
		writer.WriteByte('\n') //write endline to ensure only one per line
	}

	err = writer.Flush() //flush the final numbers still left in the buffer
	if err != nil {
		log.Printf("Final flush failed for some reason :/")
	}
	// w.reportResult()
	log.Printf("Successfully outpputted file, %d", primecount)
	w.filesToMerge = append(w.filesToMerge, outName)
	localMap = nil
	keys = nil
	runtime.GC() // for some reason it was taking up too much ram, bad for my weak PC.

	return outName, nil
}

func (w *Worker) ReportResult(success bool, resultPath string, summary string) error {
	if w.Assignment == nil {
		return fmt.Errorf("no assignment to report on")
	}

	req := &primerpc.ReportResultReq{
		WorkerID:     w.WorkerID,
		AssignmentId: w.Assignment.AssignmentID,
		ResultId:     resultPath,
		Success:      success,
		Attempt:      w.Assignment.Attempt,
	}
	resp := &primerpc.ReportResultRes{}

	log.Printf("Reporting result for task: %s, Success: %t", req.AssignmentId, req.Success)
	err := w.coordClient.Call("Coordinator.ReportResult", req, resp)
	if err != nil {

		if connShutDown(err) {
			log.Printf("conn closed trying again")
			recErr := w.reconnect()
			if recErr != nil {
				log.Printf("Could never reconnect worker %s,", w.WorkerID)
			} else {
				log.Printf("Reconnected wroker %s,", w.WorkerID)
			}

		}
	}
	if !resp.Ack {
		return fmt.Errorf("coordinator did not acknowledge")
	}
	return nil
}

/*func (w *Worker) SendHeartbeat() error {
	req := &primerpc.HeartbeatReq{
		WorkerID: w.WorkerID,
		Load:     0,
	}
	resp := &primerpc.HeartbeatRes{}

	err := w.coordClient.Call("Coordinator.Heartbeat", req, resp)
	if err != nil {
		log.Printf("Heartbeat failed: %v", err)
		return err
	}
	return nil
}

func (w *Worker) HeartbeatLoop() {
	for {
		time.Sleep(15 * time.Second)
		w.SendHeartbeat()
	}
}*/

type reader struct {
	f       *os.File
	s       *bufio.Scanner
	current uint64
}

func (w *Worker) MergeOwnFiles(bufferSize int, maxScanSize int) error {

	if len(w.filesToMerge) == 0 {
		return fmt.Errorf("Worker has no files to merge")
	}
	if len(w.filesToMerge) == 1 {
		req := &primerpc.ReportMergeReq{
			WorkerID: w.WorkerID,
			Filepath: w.filesToMerge[0],
		}
		resp := &primerpc.ReportMergeRes{}
		err := w.coordClient.Call("Coordinator.ReportMerge", req, resp)
		if err != nil {
			if connShutDown(err) {
				log.Printf("conn closed trying again")
				recErr := w.reconnect()
				if recErr != nil {
					log.Printf("Could never reconnect worker %s,", w.WorkerID)
				} else {
					log.Printf("Reconnected wroker %s,", w.WorkerID)
				}
			}
			return nil
		}

		if !resp.Ack {
			log.Printf("ReportMerge for worker %s didnt work", w.WorkerID)
		} else {
			log.Printf("ReportMerge for worker %s WORKED yaaay", w.WorkerID)
		}
		return nil
	}

	out := fmt.Sprintf("%s_%d_OwnMerge.txt", w.WorkerID, w.numOfMerges)
	if w.WorkerID == "1" {
		out = "primes1.txt"
	}
	outFile, err := w.dfsClient.Create(out)
	if err != nil {
		log.Printf("MergeFile creation failed: %v", err)
	}

	defer w.dfsClient.Close(outFile)

	writer := bufio.NewWriterSize(outFile, bufferSize)

	readers := make([]*reader, 0, len(w.filesToMerge))

	for _, filepath := range w.filesToMerge {
		file, err := w.dfsClient.Open(filepath)
		if err != nil {

			return fmt.Errorf("Could not open in worker merger %s: %v", w.WorkerID, err)
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
				return fmt.Errorf("uint parsing failed for worker %s in file %s whit err %v", w.WorkerID, i, err)
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
				return fmt.Errorf("uint parsing failed for worker %s in file %d whit err %v", w.WorkerID, item.Idx, err)
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
		w.dfsClient.Close(r.f)
	}
	log.Printf("Worker %s has merged its own files", w.WorkerID)

	w.numOfMerges++
	w.filesToMerge = w.filesToMerge[:0]

	req := &primerpc.ReportMergeReq{
		WorkerID: w.WorkerID,
		Filepath: out,
	}
	resp := &primerpc.ReportMergeRes{}
	err = w.coordClient.Call("Coordinator.ReportMerge", req, resp)
	if err != nil {
		if connShutDown(err) {
			log.Printf("conn closed trying again")
			recErr := w.reconnect()
			if recErr != nil {
				log.Printf("Could never reconnect worker %s,", w.WorkerID)
			} else {
				log.Printf("Reconnected wroker %s,", w.WorkerID)
			}
		}
		return nil
	}

	if !resp.Ack {
		log.Printf("ReportMerge for worker %s didnt work", w.WorkerID)
	} else {
		log.Printf("ReportMerge for worker %s WORKED yaaay", w.WorkerID)
	}

	return nil
}

func (w *Worker) PrimeWriteToFileSeparateFiles() (string, error) {

	assignment := w.Assignment //get assignment from the worker object
	if assignment == nil {
		return "", fmt.Errorf("No assignment available ")
	}
	filepath := assignment.TaskPath //get the filename / path

	file, err := w.dfsClient.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file to check for primes: %v", err)
	} //open the assignment file
	defer w.dfsClient.Close(file)

	localMap := make(map[uint64]bool) //create a map to ensure unique primes
	//reduces load on coordinator in the merging step

	var outFile *os.File

	scanner := bufio.NewScanner(file)
	outName := strings.TrimSuffix(filepath, fp.Ext(filepath))

	outName = fmt.Sprintf("%s_out.txt", outName)
	outFile, err = w.dfsClient.Create(outName) //create new output file
	if err != nil {
		return "", fmt.Errorf("Could not create a file to write prime output: %v", err)
	}
	defer w.dfsClient.Close(outFile)

	writer := bufio.NewWriter(outFile)
	primecount := 0 //used to track nr of primes for debugging / logs

	for scanner.Scan() {
		line := scanner.Text()
		nr, err := strconv.ParseUint(line, 10, 64) //need to convert line string to uint for prime func
		if err != nil {
			log.Printf("Failed to parse line: %q", line)
		}
		if isPrime(nr) {
			primecount++
			localMap[nr] = true

		}
	}

	keys := make([]uint64, 0, len(localMap))
	for k := range localMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, nr := range keys {
		str := strconv.FormatUint(nr, 10) //need to convert to string for writer to write it
		writer.WriteString(str)
		writer.WriteByte('\n') //write endline to ensure only one per line
	}

	err = writer.Flush() //flush the final numbers still left in the buffer
	if err != nil {
		log.Printf("Final flush failed for some reason")
	}

	log.Printf("Successfully outpputted file, %d", primecount)

	w.filesToMerge = append(w.filesToMerge, outName)

	return "", nil
}

// if it is set to true use the offset prime method, else use the separate file one
func (w *Worker) PrimeChunks(useOffsetMethod bool) (string, error) {

	if useOffsetMethod {
		return w.PrimeWriteToFile()
	} else {
		return w.PrimeWriteToFileSeparateFiles()
	}

}

func (w *Worker) ReceiveMarker(req *primerpc.SnapshotMarkerReq, resp *primerpc.SnapshotMarkerRes) error {
	w.mutex.Lock()
	if w.stateRecorded {
		resp.Ack = true
		return nil
	}

	log.Printf("Recceived marker message, w %s", w.WorkerID)

	snapshot := &primerpc.WorkerSnapshot{
		WorkerID:     w.WorkerID,
		Assignment:   w.Assignment,
		FilesToMerge: w.filesToMerge,
	}

	w.stateRecorded = true
	w.mutex.Unlock()
	go func() {
		reportResp := &primerpc.ReportResultRes{}

		err := w.coordClient.Call("Coordinator.CollectWorkerSnapshot", snapshot, reportResp)
		if err != nil {
			if connShutDown(err) {
				log.Printf("conn closed trying again")
				recErr := w.reconnect()
				if recErr != nil {
					log.Printf("Could never reconnect worker %s,", w.WorkerID)
				} else {
					log.Printf("Reconnected wroker %s,", w.WorkerID)
				}
			}
			log.Printf("Worker maybe %s failed to report snapshot: %v", w.WorkerID, err)
		} else {
			log.Printf("Worker %s successfully reported local snapshot.", w.WorkerID)
		}
	}()
	resp.Ack = true
	return nil
}

func (w *Worker) ResetSnapshotState(req *primerpc.SnapshotMarkerReq, resp *primerpc.SnapshotMarkerRes) error {
	w.mutex.Lock()
	w.stateRecorded = false
	w.mutex.Unlock()

	resp.Ack = true
	log.Printf("Worker %s state reset for next snapshot.", w.WorkerID)
	return nil
}

func connShutDown(err error) bool {
	if err == nil {
		return false
	}
	a := err.Error()
	ret := strings.Contains(a, "connection is shut down") || strings.Contains(a, "EOF")
	return ret
}

func (w *Worker) reconnect() error {
	if w.coordClient != nil {
		_ = w.coordClient.Close()
	}
	var err error
	for i := 0; i < 30; i++ {
		w.coordClient, err = rpc.Dial("tcp", w.coordAddr)
		if err == nil {
			return nil
		}
		log.Printf("Worker %s could not reconnect att %d to coord, %v", w.WorkerID, i, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("Could not reconnect at all")
}

func (w *Worker) snapshotRecovery() error {
	file, err := w.dfsClient.Open("snapshot.json")
	if err != nil {
		return fmt.Errorf("Could not open snapshot from dfs,%v", err)
	}
	defer w.dfsClient.Close(file)

	var GlobalSnapshot primerpc.GlobalSnapshot
	err1 := json.NewDecoder(file).Decode(&GlobalSnapshot)
	if err != nil {
		return fmt.Errorf("Decode dfs snapshot failes whyyy %v", err1)
	}

	localSnapshot, worked := GlobalSnapshot.WorkerStates[w.WorkerID]
	if !worked {
		log.Printf("Nahh, worker is not in snapshot")
		return nil
	}

	w.Assignment = localSnapshot.Assignment
	w.filesToMerge = localSnapshot.FilesToMerge
	return nil

}
