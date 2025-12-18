package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	worker "ds-uoe-vash/prime-app/worker"
)

func startWorker(id int, dfsAddr string, coordAddr string, basePort int) {

	workId := fmt.Sprintf("%d", id)
	port := basePort + id
	addr := fmt.Sprintf("localhost:%d", port)
	w := worker.ConstructWorker(workId, dfsAddr, addr, coordAddr)

	if w == nil {
		log.Printf("Worker %s reg failed ", workId)
		return
	}

	workerID := w.WorkerID

	// Register coord as rcp
	server := rpc.NewServer()
	err := server.RegisterName("Worker", w)

	if err != nil {
		log.Fatalf("Rpc registration in shambles %v, ", err)
	}

	listener, err := net.Listen("tcp", addr) //open port to listen
	if err != nil {
		log.Fatal("Listener error:", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Connection accept error:", err)
				continue
			}
			go server.ServeConn(conn)
		}
	}()

	log.Println("Worker RPC server listening on port, %s", addr)

	//go w.HeartbeatLoop()
	noTaskCount := 0
	for {
		err := w.RequestTask()
		if err != nil {
			log.Printf("Error requesting task: %v. Retrying after 10s", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if w.Assignment == nil {
			log.Println("No tasks available")
			noTaskCount++
			log.Printf("Worker %s: No tasks available (%d times)", workerID, noTaskCount)

			if noTaskCount >= 3 {
				log.Printf("Worker %s: Merging its own files", workerID)
				err := w.MergeOwnFiles(64*1024, 1024*64)
				if err != nil {
					log.Printf("Worker %s merging failed :'( %v", workerID, err)
				} else {
					log.Printf("Worker %s merge own files yaaay", workerID)
				}
				noTaskCount = 0
			}
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Starting task: %s", w.Assignment.TaskPath)

		// Call the modified PrimeWriteToFile which returns (string, error), if true prime write to chunks, else offsest version
		outName, err := w.PrimeChunks(false)

		if err != nil {
			log.Printf("Failed to complete task %s: %v", w.Assignment.TaskPath, err)
			reportErr := w.ReportResult(false, outName, err.Error())
			if reportErr != nil {
				log.Printf("Failed to report: %v", reportErr)

				time.Sleep(30 * time.Second)
			}
		} else {
			log.Printf("Completed task: %s -> %s", w.Assignment.TaskPath, outName)
			reportErr := w.ReportResult(true, outName, "Hurraa")
			if reportErr != nil {
				log.Printf("Failed to report  %v", reportErr)
				time.Sleep(30 * time.Second)
			}
		}
		w.Assignment = nil
	}
}

func main() {
	dfsAddr := "localhost:8090"
	coordAddr := "localhost:12345"
	basePort := 12345

	for i := 1; i <= 1; i++ {
		go startWorker(i, dfsAddr, coordAddr, basePort)
	}

	select {}
}
