// benchdata generates deterministic pseudo-random files for benchmarking.
// Given the same seed + size, it always produces identical output.
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func main() {
	seed := flag.Int64("seed", 42, "PRNG seed")
	sizeMB := flag.Int("size", 100, "output size in MiB")
	output := flag.String("out", "", "output file path (required)")
	// For partial generation: write only a portion starting at an offset.
	offsetMB := flag.Int("offset", 0, "start offset in MiB (for partial overwrite)")
	countMB := flag.Int("count", 0, "number of MiB to write (0 = full size)")
	flag.Parse()

	if *output == "" {
		fmt.Fprintln(os.Stderr, "--out is required")
		os.Exit(1)
	}

	totalSize := int64(*sizeMB) * (1 << 20)
	offset := int64(*offsetMB) * (1 << 20)
	count := int64(*countMB) * (1 << 20)
	if count == 0 {
		count = totalSize
	}

	// Create a deterministic PRNG from the seed.
	// We advance the PRNG to the offset position so that the bytes at any
	// given offset are always the same regardless of whether we're writing
	// the full file or a partial overwrite.
	rng := rand.New(rand.NewSource(*seed))

	var flags int
	if offset > 0 {
		flags = os.O_WRONLY // partial overwrite, file must exist
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(*output, flags, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open %s: %v\n", *output, err)
		os.Exit(1)
	}
	defer f.Close()

	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			fmt.Fprintf(os.Stderr, "seek: %v\n", err)
			os.Exit(1)
		}
		// Advance PRNG to offset position (in 8-byte steps since Int63 produces 8 bytes).
		steps := offset / 8
		for i := int64(0); i < steps; i++ {
			rng.Int63()
		}
	}

	buf := make([]byte, 1<<20) // 1 MiB write buffer
	written := int64(0)
	for written < count {
		// Fill buffer from PRNG.
		for i := 0; i < len(buf); i += 8 {
			v := rng.Int63()
			buf[i] = byte(v)
			buf[i+1] = byte(v >> 8)
			buf[i+2] = byte(v >> 16)
			buf[i+3] = byte(v >> 24)
			buf[i+4] = byte(v >> 32)
			buf[i+5] = byte(v >> 40)
			buf[i+6] = byte(v >> 48)
			buf[i+7] = byte(v >> 56)
		}

		toWrite := int64(len(buf))
		if written+toWrite > count {
			toWrite = count - written
		}

		n, err := f.Write(buf[:toWrite])
		if err != nil {
			fmt.Fprintf(os.Stderr, "write: %v\n", err)
			os.Exit(1)
		}
		written += int64(n)
	}

	fmt.Fprintf(os.Stderr, "wrote %d MiB to %s (seed=%d, offset=%dMiB)\n",
		written/(1<<20), *output, *seed, *offsetMB)
}
