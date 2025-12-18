## Temp main files
The main miles for both constructor and worker are temporrary, just to make sure communication works and and jobs can be assigned.

## Instructions
To properly test: go into in folder -> run this command: seq 1 100000000 > numbers.txt (it generates first 100 million numbers for testing)
1. Create dfs server
2. In new terminal create coordinator with main.go file. Just use: go run main.go (in cmd/coordinator)
3. In another new terminal create worker its main file. Use go run main.go (in cmd/worker) 
WARNING : Running the coordinator creates a new directory  cmd/coordinator/cache AND creates all the chunks in both in folder and cache folder

## Capabilities

COORDINATOR : 
a. Has two modes to chunk the files. 
    1. In a datacenter setting where the workers are in a single machine and share cache: Chunk with offsets (set numlines = 0, chunkSize = desired byte size of offset chunk)
    2. In a more distributed environment with different machines: Chunk with physical files - This places load in Fileserver and file IO (set numlines = nr of lines per file, chunkSize = 0) WARNING - This method creates many smaller files, delete the caches and the files after each use.
b. Assigns waits to be contacted by worker, assigns jobs, tracks worker merged files
c. Does the final merge primes.txt

WORKER : 
a. Processes each chunks, in the two aforementioned different modes. Set the boolean to True if processing file offsets, false otherwise
b. When it doesnt receive a new job three times in a row, it merges its own files to reduce DFS latency and load, then sends it to DFS and tells Coordinator
    (I want to fix this because if coordinator fails it still starts merging then it breakes the Coordinator merging logic)

