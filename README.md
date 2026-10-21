# go-load-balancer

A small HTTP load balancer written in Go using only the standard library. It sits in front of a set of backend servers, forwards incoming requests via reverse proxy, tracks backend health, and routes each client to a backend by hashing the client's IP address (sticky sessions / IP-hash routing).

## Features

- **Reverse proxying** — incoming requests are forwarded to backends with `httputil.NewSingleHostReverseProxy`.
- **Sticky sessions (IP hash)** — the client IP from `RemoteAddr` is hashed with FNV-32a, so requests from the same client consistently hit the same backend.
- **Failover** — if the hashed backend is dead, the pool linearly probes forward to the next alive server.
- **Health checking** — every backend is polled at `<backend URL>/health` every 20 seconds (1-second timeout, expecting `200 OK`).
- **Failure detection** — after 3 consecutive proxy errors a backend is marked dead; any successful response resets the counter.
- **Zero dependencies** — standard library only.

## Quick start

Start two sample backends (each exposes `GET /health` and a root handler that identifies itself):

```sh
go run ./testserver -port 8001
go run ./testserver -port 8002
```

Start the load balancer, pointing it at the backends:

```sh
go run . -servers http://localhost:8001,http://localhost:8002
```

The load balancer now accepts traffic on `http://localhost:8000`:

```sh
curl http://localhost:8000
# response from backend-8001 (path: /)
```

Backends are expected to expose a `GET /health` endpoint returning HTTP 200.

## Verifying sticky sessions and failover

The `testuser` tool simulates multiple users with distinct source IPs (bound via `net.Dialer.LocalAddr`) and prints a per-user breakdown of which backend answered:

```sh
go run ./testuser -n 10 -ips 127.0.0.1,127.0.0.2
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-target` | `http://localhost:8000` | Load balancer address |
| `-n` | `100` | Requests per user (one per second) |
| `-ips` | *(one default client)* | Comma-separated source IPs to bind to |

Each source IP should consistently hit the same backend. Kill one of the `testserver` processes and re-run — its traffic should fail over to the surviving backend, and the health checker should mark it dead within 20 seconds.

The `testserver` tool flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8003` | Port to listen on (`127.0.0.1:<port>`) |
| `-name` | `backend-<port>` | Name reported in responses |

## Running the tests

```sh
go test ./...              # all tests
go test ./loadbalancer/ -v # verbose, per-test output
go test ./... -race        # with the race detector
```

The unit tests cover server-pool routing (hash mapping, dead-server probing, wrap-around), server state transitions (failure threshold, health checks against `httptest` backends), and request processing (503 on empty pool, end-to-end proxying).

## Project structure

```
main.go                      CLI entry point — parses -servers, starts the balancer
loadbalancer/
    load-balancer.go         LoadBalancer: listen address, IP hashing, request dispatch
    serverpool.go            ServerPool: modulo-hash server selection with failover probing
    server.go                Server: reverse proxy, alive flag, failure counting, health checks
testserver/main.go           Sample backend for manual testing
testuser/main.go             Load-testing client simulating users with distinct source IPs
```

## Configuration constants

Defined in `loadbalancer/`:

| Constant | Value | Meaning |
|----------|-------|---------|
| `LOAD_BALANCER_ADDRESS` | `localhost:8000` | Address the balancer listens on |
| `MAX_FAILURES` | `3` | Consecutive proxy errors before a backend is marked dead |
