package main

import (
	"fmt"
	"log"
	"os"

	"ds-uoe-vash/dfs/pkg/client"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <cache-directory> <server-address> [server-address...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s /tmp/afs localhost:8080 localhost:8081 localhost:8082\n", os.Args[0])
		os.Exit(1)
	}

	cacheDir := os.Args[1]
	serverAddrs := os.Args[2:]

	dfsClient, err := client.NewDFSClient(serverAddrs, cacheDir)
	if err != nil {
		log.Fatalf("Failed to initialize DFS client: %v", err)
	}
	defer dfsClient.Disconnect()

	log.Printf("DFS Client initialized successfully")
	log.Printf("  Server Addresses: %v", serverAddrs)
	log.Printf("  Cache Directory:  %s", cacheDir)
	log.Printf("  Ready to use DFS operations")

	runTests(dfsClient)
}

func runTests(dfsClient *client.DFSClient) {
	fmt.Println("\n=== DFS Client Tests ===\n")

	fmt.Println("1. Opening a file from server...")
	file, err := dfsClient.Open("input.txt")
	if err != nil {
		log.Printf("   Could not open file: %v", err)
	} else {

		buf := make([]byte, 1024)
		n, _ := dfsClient.Read(file, buf)
		fmt.Printf("   Read %d bytes from cached file\n", n)

		dfsClient.Close(file)
		fmt.Println("   File closed")
	}

	fmt.Println("\n2. Creating a new output file...")
	newFile, err := dfsClient.Create("output.txt")
	if err != nil {
		log.Printf("   Could not create file: %v", err)
	} else {

		data := []byte("Hello from DFS client!\nThis is cached locallly.\n")
		n, _ := dfsClient.Write(newFile, data)
		fmt.Printf("   Wrote %d bytes to cached file\n", n)

		fmt.Println("   Closing file (flushing to server)...")
		if err := dfsClient.Close(newFile); err != nil {
			log.Printf("   Error during flush: %v", err)
		} else {
			fmt.Println("   File successfully flushed to server")
		}
	}

	fmt.Println("\n3. Reopening the same file (should use cache)...")
	cachedFile, err := dfsClient.Open("input.txt")
	if err != nil {
		log.Printf("   Could not open file: %v", err)
	} else {
		fmt.Println("   File opened from cache")
		dfsClient.Close(cachedFile)
	}

	fmt.Println("\n=== Demonstration Complete ===")
	fmt.Println("\nYou can now integrate this client with your application.")
	fmt.Println("Available operations:")
	fmt.Println("  - dfsClient.Open(filename)")
	fmt.Println("  - dfsClient.Create(filename)")
	fmt.Println("  - dfsClient.Read(file, buffer)")
	fmt.Println("  - dfsClient.Write(file, data)")
	fmt.Println("  - dfsClient.Close(file)")
	fmt.Println("  - dfsClient.ReadAt(file, buffer, offset)")
	fmt.Println("  - dfsClient.WriteAt(file, data, offset)")
	fmt.Println("  - dfsClient.Seek(file, offset, whence)")
}
