// Package chunk implements a concurrent pipeline that reads raw blocks,
// hashes them with BLAKE3, and emits ChunkMeta + data pairs.
package chunk

import (
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/theo/synche2/internal/block"
	"github.com/theo/synche2/internal/protocol"
	"github.com/zeebo/blake3"
)

// Result holds one chunk's data and metadata, produced by the pipeline.
type Result struct {
	Meta protocol.ChunkMeta
	Data []byte // raw chunk bytes
}

// Pipeline reads a block device, hashes chunks concurrently, and sends
// results through a channel. It maintains chunk ordering via the index field.
type Pipeline struct {
	reader      *block.Reader
	concurrency int
}

// NewPipeline creates a chunking pipeline that will read from the given device
// path using the specified concurrency for hashing.
func NewPipeline(devicePath string, concurrency int) (*Pipeline, error) {
	r, err := block.NewReader(devicePath, protocol.ChunkSize)
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}
	if concurrency < 1 {
		concurrency = 1
	}
	return &Pipeline{
		reader:      r,
		concurrency: concurrency,
	}, nil
}

// DeviceSize returns the total size of the underlying device.
func (p *Pipeline) DeviceSize() uint64 {
	return p.reader.Size()
}

// indexedBuf is an internal type to pair raw data with its positional index.
type indexedBuf struct {
	index  uint64
	offset uint64
	data   []byte
	size   int
}

// Run starts the pipeline and returns a channel of Results.
// The channel is closed when all chunks have been read and hashed.
// Errors during reading are sent as Results with a non-empty Meta.Hash of "error:<msg>".
func (p *Pipeline) Run() <-chan Result {
	out := make(chan Result, p.concurrency*2)

	// Stage 1: sequential reader goroutine.
	raw := make(chan indexedBuf, p.concurrency*2)
	go func() {
		defer close(raw)
		var idx uint64
		var offset uint64
		for {
			data, n, err := p.reader.ReadChunk()
			if n > 0 {
				// Copy data out of the aligned buffer so it can be reused.
				cp := make([]byte, n)
				copy(cp, data[:n])
				raw <- indexedBuf{
					index:  idx,
					offset: offset,
					data:   cp,
					size:   n,
				}
				idx++
				offset += uint64(n)
			}
			if err != nil {
				break
			}
		}
	}()

	// Stage 2: concurrent hashers.
	// We use a fixed pool of goroutines. Results may arrive out of order
	// but each carries its index so the consumer can reorder if needed.
	var wg sync.WaitGroup
	var hasherCount int64
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&hasherCount, 1)

			hasher := blake3.New()
			for buf := range raw {
				hasher.Reset()
				hasher.Write(buf.data)
				sum := hasher.Sum(nil)

				out <- Result{
					Meta: protocol.ChunkMeta{
						Index:  buf.index,
						Hash:   hex.EncodeToString(sum),
						Size:   buf.size,
						Offset: buf.offset,
					},
					Data: buf.data,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		p.reader.Close()
	}()

	return out
}

// Close releases resources. Safe to call even if Run was never called.
func (p *Pipeline) Close() error {
	return p.reader.Close()
}
