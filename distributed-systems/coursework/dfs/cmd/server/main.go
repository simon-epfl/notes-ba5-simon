package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"ds-uoe-vash/dfs/pkg/server"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-directory> <output-directory> [port] [peer-addrs...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ./input ./output 8080 localhost:8081 localhost:8082\n", os.Args[0])
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]
	port := "8090"
	var peerAddrs []string

	if len(os.Args) >= 4 {
		port = os.Args[3]
	}

	if len(os.Args) >= 5 {
		peerAddrs = os.Args[4:]
	}

	if err := validateDirectory(inputDir); err != nil {
		log.Fatalf("Input directory error: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fileServer := server.NewFileServer(inputDir, outputDir, peerAddrs)

	if err := rpc.Register(fileServer); err != nil {
		log.Fatalf("Failed to register RPC service: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("DFS Server started successfully")
	log.Printf("  Input Directory:  %s", inputDir)
	log.Printf("  Output Directory: %s", outputDir)
	log.Printf("  Listening on:     :%s", port)
	if len(peerAddrs) > 0 {
		log.Printf("  Peer servers:     %v", peerAddrs)
	}
	log.Printf("  Server ready to accept connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func validateDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}
