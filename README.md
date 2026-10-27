# Go Profiling with pprof

Code for the video [Golang Profiling with pprof: CPU Flame Graphs and Memory Leaks](https://youtu.be/As8Jdq9Gpwg).

A log search service (`/search?q=<regexp>` over 20,000 generated log lines) with two bugs that pprof finds:

- `server/`: compiles the regexp inside the loop (CPU) and keeps every response in the `recent` slice (memory leak)
- `fixed/`: compiles once before the loop, returns 400 for a bad pattern, and keeps only the last 50 results

Both import `net/http/pprof`, so the profiles are at `localhost:9092/debug/pprof/`.

## Run it in a Linux container

On macOS the CPU profile is full of runtime calls from the darwin port, so the video runs the service in a container.

```bash
GOOS=linux go build -o bin/server-linux ./server
GOOS=linux go build -o bin/fixed-linux ./fixed
docker build -t goprof:v1 --build-arg BIN=server .
docker build -t goprof:v2 --build-arg BIN=fixed .
docker run -d --name goprof -p 9092:9092 --cpus 4 -m 1g goprof:v1
```

## CPU profile

```bash
k6 run -q k6/search.js        # 4 virtual users for 20 s
go tool pprof -http=:9093 "localhost:9092/debug/pprof/profile?seconds=10"
```

Run the second command while k6 is running. In the flame graph, `regexp.MustCompile` takes most of `main.search`.

## Heap profile

```bash
k6 run -q k6/steady.js        # 10 searches a second for 60 s
docker stats goprof           # memory climbs with server/, stays flat with fixed/
curl -s localhost:9092/debug/pprof/heap > heap1.pb.gz
# a minute later
curl -s localhost:9092/debug/pprof/heap > heap2.pb.gz
go tool pprof -http=:9094 -base heap1.pb.gz heap2.pb.gz
```

With `-base`, pprof subtracts the first profile, so only the memory added in between is shown.

Swap `goprof:v1` for `goprof:v2` to run the same tests against the fixed service.

The `/debug/pprof` handlers sit on the default mux here. Don't do that on a public server: give pprof its own port that only your network can reach.
