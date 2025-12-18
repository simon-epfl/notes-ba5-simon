package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"time"

	coordinator "ds-uoe-vash/prime-app/coordinator"
)

func triggerSnapshot(c *coordinator.Coordinator) {
	const interval = 5 * time.Second
	for {
		time.Sleep(interval)
		err := c.StartSnapshot()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func main() {

	var linesPerFile = 100000
	var chunksize = 0
	var numfiles = 1
	var start = 1

	// ip and root and cache is currently hardcoded, add later as cmd argument
	c := coordinator.ConstructCoordinator("localhost:8090", "/dfs-root", 10)

	// register the Coordinator as an RPC service
	err := rpc.Register(c)
	if err != nil {
		log.Fatal("Error registering RPC:", err)
	}
	go triggerSnapshot(c)
	// listen on a TCP port
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal("Listener error:", err)
	}
	defer listener.Close()

	log.Println("Coordinator RPC server listening on port 12345")

	if !c.Recovered {
		log.Println("STarted chunking file")

		for i := start; i <= numfiles; i++ { //set back to 1
			filename := fmt.Sprintf("input_dataset_%03d.txt", i)
			log.Printf("Chunking %s", filename)

			err = c.ChunkFiles(filename, linesPerFile, chunksize)
			if err != nil {
				log.Printf("ERROR: Failed to chunk file: %v", err)
			} else {
				log.Println("File chunked successfully")
			}
		}
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

}
