# synche

A concurrent, deduplicated file upload system written in Go. Reads raw data from files or block devices, chunks it into 1 MiB pieces, hashes each chunk with BLAKE3, and uploads only novel chunks to a server. Uploaded files can be browsed and downloaded via WebDAV.

## Why

- **Faster than rsync** — concurrent hashing, probing, and uploading keeps the pipeline saturated
- **Deduplication** — identical chunks are stored once on the server, even across different files
- **Local caching** — unchanged chunks are skipped entirely on subsequent runs (no network I/O)
- **Content-addressable storage** — chunks are keyed by BLAKE3 hash, so two files that share data (e.g. a short video and a longer version of it) share chunks on disk
- **WebDAV access** — uploaded files appear as downloadable files at `/webdav/`, reassembled from chunks on the fly
- **Raw block device support** — reads directly with `O_DIRECT` to bypass the kernel page cache

## Benchmark

Tested against a remote server (YOUR_SERVER_IP), 100 MiB of deterministic data, averaged over 3 rounds:

```
                                           synche    rsync(ssh)
                                         --------  ------------
Case 1: 100 MiB fresh upload (avg)        3409 ms       4592 ms
Case 2: 100 MiB, 50 MiB exists (avg)      2293 ms       4014 ms

Case 1 vs rsync(ssh): synche is 1.3x faster
Case 2 vs rsync(ssh): synche is 1.7x faster

Case 1 individual rounds (ms):
  round 1:  synche=3425  ssh=4392
  round 2:  synche=3537  ssh=4763
  round 3:  synche=3265  ssh=4623

Case 2 individual rounds (ms):
  round 1:  synche=2447  ssh=5067
  round 2:  synche=2343  ssh=3758
  round 3:  synche=2089  ssh=3219
```

The advantage grows in Case 2 because synche skips chunks the server already has, while rsync still needs to compute and compare rolling checksums over the full file.

## Usage

### Server

```sh
go build -o bin/synche-server ./cmd/synche-server
./bin/synche-server --addr :8420 --store ./synche-store
```

Uploaded files are browsable at `http://localhost:8420/webdav/`.

### Client

```sh
go build -o bin/synche-client ./cmd/synche-client
./bin/synche-client --server http://localhost:8420 --source /path/to/file
```

Options:
- `--server` — server URL (default `http://localhost:8420`)
- `--source` — file or block device to upload
- `--concurrency` — parallel upload workers (default 20)
- `--no-cache` — disable local manifest cache

### Benchmark

```sh
./benchmark.sh --server-host <ip> --rounds 3
```

## How it works

1. **Read** — the client reads the source file in 1 MiB chunks using `O_DIRECT` (falls back to normal I/O when not supported)
2. **Hash** — chunks are hashed concurrently with BLAKE3
3. **Cache check** — each chunk hash is compared against a local cache from the previous run; unchanged chunks are skipped with zero network I/O
4. **Probe** — remaining hashes are sent to the server in batches; the server returns which ones it doesn't have
5. **Upload** — only novel chunks are uploaded, with retries on failure
6. **Manifest** — a manifest mapping chunk indices to hashes is saved on the server

The server stores chunks in a content-addressable filesystem (`chunks/ab/cd/<hash>`). Multiple manifests can reference the same chunks, so uploading similar files costs only the storage of their unique chunks.

## Project structure

```
cmd/
  synche-client/    Client CLI
  synche-server/    Server CLI
  benchdata/        Deterministic test data generator
internal/
  block/            Raw block device I/O with O_DIRECT
  cache/            Local manifest cache (~/.cache/synche/)
  chunk/            Concurrent read → hash pipeline
  client/           Upload client with probe, retry, dedup
  protocol/         Shared types and constants
  server/           HTTP API server
  store/            Content-addressable chunk storage
  webdav/           Read-only WebDAV handler
```
