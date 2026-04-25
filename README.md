# redis-go

> A protocol-compliant Redis server written from scratch in Go — zero external dependencies, pure stdlib, generic type-safe store, AOF persistence, and full container/K8s deployment.

```
$ redis-cli PING
PONG

$ redis-cli XADD stream_key 1526919030474-0 temperature 36 humidity 95
"1526919030474-0"

$ redis-cli BLPOP queue 2.5
1) "queue"
2) "job-payload"
```

---

## What Makes This Different

Most Redis reimplementations stop at `PING`, `SET`, and `GET`. This one goes further — and the interesting parts aren't the commands, they're the decisions underneath them.

**Single generic store, runtime type dispatch.**
Rather than maintaining separate maps for strings, lists, and streams, the entire server runs on one `Storage[any]` instance. Type identity is recovered through Go's runtime type assertions at read time. Each command handler knows exactly what type it expects — `TYPE` introspects on the concrete Go type, and type-specific operations like `XADD` or `LPUSH` assert and operate on the correct underlying structure. No union types, no separate maps, no duplication.

```go
type Storage[T any] struct {
    mu        sync.RWMutex
    store     map[string]Value[T]
    notifiers map[string][]chan struct{}  // per-key blocking subscriber queue
    Config    Config
}

var Cache = &Storage[any]{...}
```

**`BLPOP` via channel-based notification — no polling.**
When a `BLPOP` arrives and the list is empty, the handler registers a `chan struct{}` bell on the key, releases the store lock, and blocks in a `select`. When any `LPUSH`/`RPUSH` appends to that key, it rings the first waiting bell. The goroutine wakes, pops immediately, and returns — zero polling, zero busy-waiting. Timeout is a `time.After` in the same `select`.

**AOF replay goes through the live router.**
Rather than writing a separate crash-recovery parser, the startup replayer feeds AOF commands through the same `router.HandleCommand` used for live traffic — with `io.Discard` as the writer. One code path for both live and replay. The behaviour is always identical.

**Expiry without a background sweeper.**
TTLs are stored as absolute Unix-millisecond deadlines computed at write time. `GET` checks the deadline on every read. No background goroutine, no timer heap, no GC pressure from periodic sweeps.

**No external dependencies.** The `go.mod` has a single module declaration. RESP parsing, TCP handling, concurrency, file I/O, and the full AOF implementation are all standard library.

---

## Architecture

```
app/
├── main.go                        # TCP listener, goroutine-per-connection loop
├── lib/
│   ├── lib.go                     # RESP marshal / unmarshal
│   └── utils.go                   # Stream ID parsing and range comparison
├── lib/commands/
│   ├── router/handler.go          # Command dispatch router
│   ├── set.go, get.go, ping.go …  # Per-command handlers
│   └── xadd.go, xrange.go …       # Stream command handlers
└── store/
    ├── store.go                   # Generic Storage[T] definition and init
    └── utils.go                   # All storage operations + AOF persistence
```

One goroutine per TCP connection. A single `sync.RWMutex` guards all store access — readers share a read lock, writers take an exclusive lock. Blocking operations release before entering `select` to avoid deadlocking against concurrent pushers.

---

## Features

### Core Protocol

- Full **RESP (REdis Serialization Protocol)** parser — arrays, bulk strings, inline commands, nested structures
- Custom `MarshalArrayRESP` / `UnmarshalRESP`, including recursive marshalling of nested arrays
- Buffered reader per connection for stream-safe incremental parsing

### String Commands

| Command | Notes                                             |
| ------- | ------------------------------------------------- |
| `PING`  |                                                   |
| `ECHO`  |                                                   |
| `SET`   | `PX` (millisecond) and `MX` (second) expiry flags |
| `GET`   | Deadline-aware expiry check on every read         |

### List Commands

| Command           | Notes                                                    |
| ----------------- | -------------------------------------------------------- |
| `LPUSH` / `RPUSH` | Multi-element push, correct ordering semantics           |
| `LPOP`            | Single or count-based pop                                |
| `LRANGE`          | Full negative-index support and Redis boundary semantics |
| `LLEN`            |                                                          |
| `BLPOP`           | Blocking pop with timeout via channel notification       |

### Stream Commands

| Command  | Notes                                                                |
| -------- | -------------------------------------------------------------------- |
| `XADD`   | Full ID validation, partial auto-ID (`ms-*`), full auto-ID (`*`)     |
| `XRANGE` | Range queries with `-` / `+` sentinel support                        |
| `TYPE`   | Returns `string`, `list`, or `stream` via runtime type introspection |

Stream entry IDs follow the `<milliseconds>-<sequence>` format. Auto-IDs are resolved at write time by reading the last entry's timestamp and incrementing the sequence monotonically — the server enforces the invariant that IDs are always strictly increasing.

### Append-Only File (AOF) Persistence

The persistence layer mirrors real Redis AOF architecture:

- An **incremental `.incr.aof` file** captures every write command in RESP format
- A **manifest file** tracks the currently active `.incr.aof` file
- On startup, the server **replays** the AOF file through the live command router — no separate replay parser
- `appendfsync always` forces an OS `fsync` after every write; `everysec` defers to the page cache
- All config is runtime-injectable via CLI flags

```
data/backup/
├── backup.aof.manifest          # Points to the active incremental file
└── backup.aof.1.incr.aof        # RESP-encoded write-ahead log
```

### CONFIG

`CONFIG GET` supported for: `dir`, `appendonly`, `appenddirname`, `appendfilename`, `appendfsync`, `dbfilename`.

---

## Running Locally

**Requirements:** Go 1.21+

```bash
git clone https://github.com/0xjustuzair/redis-go
cd redis-go

# Without persistence
./your_program.sh

# With AOF persistence
./your_program.sh \
  --dir data \
  --appendonly yes \
  --appenddirname backup \
  --appendfilename backup.aof \
  --appendfsync always

# Connect
redis-cli PING
redis-cli SET hello world
redis-cli GET hello
redis-cli XADD mystream '*' sensor temperature
```

### Makefile

```bash
make run          # Start server
make ping         # redis-cli PING
make set          # SET foo bar
make blpop        # RPUSH + BLPOP
make xadd         # XADD stream_key ...
make xrange       # XADD x2 + XRANGE
make persis-run   # Start with AOF persistence
make test-all     # Run all manual test targets
```

---

## Deployment

### Docker

```bash
docker build -t redis-go .
docker run -p 6379:6379 redis-go
```

### Docker Compose

```bash
docker-compose up
```

### Kubernetes (StatefulSet)

Deployed as a `StatefulSet` — not a `Deployment` — because Redis is stateful: each pod needs a stable identity and a persistent volume claim for AOF durability. The manifest provisions a 1Gi PVC, configures `tcpSocket` readiness probes on port 6379, and sets conservative resource requests (128Mi memory, 0.1 CPU) appropriate for a Go process.

```bash
kubectl apply -f app/k8s/deployment.yaml
kubectl apply -f app/k8s/service.yaml
```

---

## License

MIT
