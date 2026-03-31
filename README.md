# synche

A concurrent, deduplicated file upload system written in Go. Reads raw data from files or block devices, chunks it into 1 MiB pieces, hashes each chunk with BLAKE3, and uploads only novel chunks to a server. Uploaded files can be browsed and downloaded via WebDAV.

## Why

- **Faster than rsync** — concurrent hashing, probing, and uploading keeps the pipeline saturated
- **Deduplication** — identical chunks are stored once on the server, even across different files
- **Local caching** — unchanged chunks are skipped entirely on subsequent runs (no network I/O)
- **Content-addressable storage** — chunks are keyed by BLAKE3 hash, so two files that share data (e.g. a short video and a longer version of it) share chunks on disk
- **WebDAV access** — uploaded files appear as downloadable files at `/webdav/`, reassembled from chunks on the fly
- **Resumable uploads** — if an upload is interrupted, re-running the same command picks up where it left off; chunks already on the server are skipped automatically
- **Raw block device support** — reads directly with `O_DIRECT` to bypass the kernel page cache

## Benchmark

Tested against a remote server, 100 MiB of deterministic data, averaged over 3 rounds:

```
                                               synche      rsync(ssh)             scp
                                             --------  --------------  --------------
  Case 1: 100 MiB fresh upload (avg)          8699 ms        16558 ms        13004 ms
  Case 2: 100 MiB, 50 MiB exists (avg)        4655 ms         6401 ms        18660 ms

  Case 1 vs rsync(ssh): synche is 1.9x faster
  Case 2 vs rsync(ssh): synche is 1.3x faster
  Case 1 vs scp: synche is 1.4x faster
  Case 2 vs scp: synche is 4.0x faster

  Case 1 individual rounds (ms):
    round 1:  synche=8703  ssh=17017  scp=9909
    round 2:  synche=8720  ssh=14956  scp=18808
    round 3:  synche=8676  ssh=17701  scp=10295

  Case 2 individual rounds (ms):
    round 1:  synche=4569  ssh=6347  scp=35562
    round 2:  synche=4969  ssh=7057  scp=10437
    round 3:  synche=4429  ssh=5800  scp=9981
```

scp is a raw SSH pipe with no checksumming or delta logic — it's the baseline for "how fast can bytes travel over SSH." synche beats it in both cases because it skips known chunks. The advantage is dramatic in Case 2: scp always sends the full file, while synche skips the 50 MiB already on the server.

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
- `--source` — file or block device to upload
- `--api-key` — API key for server authentication (or set `SYNCHE_API_KEY` env var)
- `--concurrency` — parallel upload workers (default: number of CPUs)
- `--no-cache` — disable local manifest cache

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
