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
Case 1: 100 MiB fresh upload (avg)          3925 ms         4806 ms         4798 ms
Case 2: 100 MiB, 50 MiB exists (avg)        2640 ms         3853 ms         5817 ms

Case 1 vs rsync(ssh): synche is 1.2x faster
Case 2 vs rsync(ssh): synche is 1.4x faster
Case 1 vs scp: synche is 1.2x faster
Case 2 vs scp: synche is 2.2x faster

Case 1 individual rounds (ms):
  round 1:  synche=3477  ssh=4703  scp=4405
  round 2:  synche=3926  ssh=4567  scp=4867
  round 3:  synche=4372  ssh=5149  scp=5124

Case 2 individual rounds (ms):
  round 1:  synche=2399  ssh=4037  scp=5484
  round 2:  synche=2836  ssh=4038  scp=5359
  round 3:  synche=2687  ssh=3484  scp=6609
```

scp is a raw SSH pipe with no checksumming or delta logic — it's the baseline for "how fast can bytes travel over SSH." synche beats it in both cases because it skips known chunks. The advantage is dramatic in Case 2: scp always sends the full file, while synche skips the 50 MiB already on the server.

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
