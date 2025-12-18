## Installation

```go
import "ds-uoe-vash/dfs/pkg/client"
```

## Usage

```go
package main

import (
    "log"
    "ds-uoe-vash/dfs/pkg/client"
)

func main() {
    // Connect to server
    dfsClient, err := client.NewDFSClient("localhost:8080", "./cache")
    if err != nil {
        log.Fatal(err)
    }
    defer dfsClient.Disconnect()

    // Open existing file (fetches from server, caches locally)
    file, err := dfsClient.Open("example.txt")
    if err != nil {
        log.Fatal(err)
    }

    // Read from local cache (no requests to the server here!)
    buf := make([]byte, 1024)
    n, _ := dfsClient.Read(file, buf)
    log.Printf("Read %d bytes: %s", n, buf[:n])

    // Write to local cache (no requests to the server here!)
    dfsClient.Write(file, []byte("\nNew line added"))

    // Close - flushes changes to server (here you send it to the server!)
    if err := dfsClient.Close(file); err != nil {
        log.Fatal(err)
    }
}
```

## API

```go
// Core operations
NewDFSClient(serverAddr, cacheDir)  // Connect to server
Open(filename)                       // Open existing file
Create(filename)                     // Create new file
Read(file, buffer)                   // Read from cache
Write(file, data)                    // Write to cache
Close(file)                          // Flush and close
Disconnect()                         // Close connection

// Additional operations
ReadAt(file, buf, offset)            // Random access read
WriteAt(file, data, offset)          // Random access write
Seek(file, offset, whence)           // Move cursor
ReadFull(file)                       // Read entire file
WriteAll(file, data)                 // Replace entire content
```

## How It Works

1. **Open**: Fetches entire file from server, caches locally
2. **Read/Write**: All operations on local cache (fast, no network)
3. **Close**: Sends modified files back to server
4. **Retry**: Automatically retries failed requests up to 3 times

Files must always be opened with the client before you start working.
The client will retry 3 times every 100ms if the server doesn't respond