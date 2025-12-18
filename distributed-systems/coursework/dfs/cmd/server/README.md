## Usage

```bash
go run main.go <input-directory> <output-directory> [port]
```

- **input-directory**: Where the server reads files from
- **output-directory**: Where the server writes modified files
- **port**: Optional, defaults to 8080

## Example

```bash
# Create directories
mkdir -p input output

# Add a test file
echo "Hello, DFS!" > input/test.txt

# Start server
go run main.go ./input ./output 8080
```

The server will:
- Serve files from `./input`
- Write client changes to `./output`
- Listen on port 8080 for RPC connections

## Testing with a Client

```go
// Connect
client, _ := client.NewDFSClient("localhost:8080", "./cache")

// Open file from input directory
file, _ := client.Open("test.txt")

// Read and modify
buf := make([]byte, 100)
client.Read(file, buf)
client.Write(file, []byte("Hello Arnold and Haroon")) <- local, only change the file on the client

// Close - writes to output directory       
client.Close(file)                                    <- flushes the new content on the server
```

After closing, check `./output/test.txt` for the changes.
Remember that the changes are done in local first and then propagated to the server once you close the file
