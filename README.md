# synche

A concurrent, deduplicated file upload system written in Go. Reads raw data from files or block devices, chunks it into 1 MiB pieces, hashes each chunk with BLAKE3, and uploads only novel chunks to a server. Uploaded files can be browsed and downloaded via WebDAV.

## Why

- **Faster than rsync** — concurrent hashing, probing, and uploading keeps the pipeline saturated
- **Deduplication** — identical chunks are stored once on the server, even across different files
- **Local caching** — unchanged chunks are skipped entirely on subsequent runs (no network I/O)
- **Content-addressable storage** — chunks are keyed by BLAKE3 hash, so two files that share data (e.g. a short video and a longer version of it) share chunks on disk
- **WebDAV access** — uploaded files appear as downloadable files at `/webdav/`, reassembled from chunks on the fly
- **Resumable uploads** — if an upload is interrupted, re-running the same command picks up where it left off; chunks already on the server are skipped automatically
- **Directory uploads** — walks a directory tree and uploads each file individually with a 4-phase concurrent pipeline; each file gets its own manifest visible in WebDAV
- **HTTP/2** — client and server use HTTP/2 (h2c for cleartext, ALPN for TLS) to multiplex all requests over a single TCP connection
- **Raw block device support** — reads directly with `O_DIRECT` to bypass the kernel page cache

## Benchmark

Tested against a remote server over HTTP/2, averaged over 3 rounds. Six test cases covering single files, disk images, and directory uploads:

| Case | synche | rsync (ssh) | synche advantage |
|------|--------|-------------|------------------|
| 100 MiB fresh upload | 4,068 ms | 4,588 ms | **1.1x faster** |
| 100 MiB, 50 MiB already on server | 2,159 ms | 2,909 ms | **1.3x faster** |
| 90 MiB disk image fresh | 3,463 ms | 4,219 ms | **1.2x faster** |
| Disk image, 10% changed | 952 ms | 1,888 ms | **1.9x faster** |
| 100 files (directory) fresh | 3,818 ms | 4,116 ms | **1.0x faster** |
| 100 files, 10% changed | 992 ms | 1,872 ms | **1.8x faster** |

The biggest wins are on incremental syncs (Cases 4 and 6) — synche's content-addressable storage means only changed chunks need to travel over the wire. Directory uploads use a 4-phase pipeline (hash all → batch probe → upload needed chunks → upload manifests) to minimize round-trips.

## Usage

### Server

```sh
go build -o bin/synche-server ./cmd/synche-server
./bin/synche-server --addr :8420 --store ./synche-store --api-key YOUR_SECRET_KEY
```

Options:
- `--addr` — listen address (default `:8420`)
- `--store` — chunk store directory (default `./synche-store`)
- `--api-key` — require API key for all requests (or set `SYNCHE_API_KEY` env var)
- `--tls-cert` — path to TLS certificate file (enables HTTPS)
- `--tls-key` — path to TLS private key file

With TLS:
```sh
./bin/synche-server --addr :8420 --store ./synche-store \
    --api-key YOUR_SECRET_KEY \
    --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem
```

Uploaded files are browsable at `http://localhost:8420/webdav/` (or `https://` with TLS).

### Client

```sh
go build -o bin/synche-client ./cmd/synche-client
./bin/synche-client --server http://localhost:8420 --source /path/to/file --api-key YOUR_SECRET_KEY
```

Options:
- `--server` — server URL (default `http://localhost:8420`)
- `--source` — file, directory, or block device to upload
- `--api-key` — API key for server authentication (or set `SYNCHE_API_KEY` env var)
- `--concurrency` — parallel upload workers (default: number of CPUs)
- `--no-cache` — disable local manifest cache

Upload a directory:
```sh
./bin/synche-client --server http://localhost:8420 --source /path/to/dir --api-key YOUR_SECRET_KEY
```

Each file in the directory gets its own manifest and appears individually in WebDAV.

### Benchmark

```sh
./benchmark.sh --server YOUR_SERVER_IP --rsync-mode ssh --ssh-user root --api-key YOUR_SECRET_KEY
```

## How it works

1. **Read** — the client reads the source file in 1 MiB chunks using `O_DIRECT` (falls back to normal I/O when not supported)
2. **Hash** — chunks are hashed concurrently with BLAKE3
3. **Cache check** — each chunk hash is compared against a local cache from the previous run; unchanged chunks are skipped with zero network I/O
4. **Probe** — remaining hashes are sent to the server in batches; the server returns which ones it doesn't have
5. **Upload** — only novel chunks are uploaded, with retries on failure
6. **Manifest** — a manifest mapping chunk indices to hashes is saved on the server

The server stores chunks in a content-addressable filesystem (`chunks/ab/cd/<hash>`). Multiple manifests can reference the same chunks, so uploading similar files costs only the storage of their unique chunks.
