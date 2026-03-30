// Package block provides raw block device I/O using direct I/O (O_DIRECT)
// to bypass the kernel page cache and read raw disk data as fast as possible.
package block

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// alignment required for O_DIRECT buffers (typically 512 or 4096).
	alignment = 4096
)

// Reader reads raw bytes from a block device or file using O_DIRECT,
// bypassing the kernel page cache entirely.
type Reader struct {
	fd       int
	size     uint64 // total device/file size in bytes
	pos      uint64
	readSize int // how many bytes to read per call
}

// NewReader opens the given path with O_DIRECT | O_RDONLY.
// It works on block devices (/dev/sdX, /dev/nvmeXnY) and regular files.
func NewReader(path string, readSize int) (*Reader, error) {
	// Align readSize up to alignment boundary (O_DIRECT requirement).
	if readSize%alignment != 0 {
		readSize = ((readSize / alignment) + 1) * alignment
	}

	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_DIRECT, 0)
	if err != nil {
		// Fall back to non-direct I/O (useful for testing with regular files
		// on filesystems that don't support O_DIRECT, like tmpfs).
		fd, err = unix.Open(path, unix.O_RDONLY, 0)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}
	}

	size, err := deviceSize(fd, path)
	if err != nil {
		unix.Close(fd)
		return nil, err
	}

	return &Reader{
		fd:       fd,
		size:     size,
		readSize: readSize,
	}, nil
}

// Size returns the total size of the device/file in bytes.
func (r *Reader) Size() uint64 {
	return r.size
}

// ReadChunk reads the next readSize bytes into an aligned buffer.
// Returns the buffer, actual bytes read, and any error.
// At EOF, returns io.EOF with n=0.
func (r *Reader) ReadChunk() ([]byte, int, error) {
	if r.pos >= r.size {
		return nil, 0, io.EOF
	}

	buf := alignedAlloc(r.readSize)

	remaining := r.size - r.pos
	toRead := uint64(r.readSize)
	if toRead > remaining {
		// For the last chunk on O_DIRECT, we still need to read an aligned
		// amount, then truncate the result.
		toRead = remaining
	}

	// O_DIRECT requires aligned read sizes. Round up for the syscall.
	alignedRead := toRead
	if alignedRead%uint64(alignment) != 0 {
		alignedRead = ((alignedRead / uint64(alignment)) + 1) * uint64(alignment)
	}
	// Make sure we don't read past buf capacity.
	if alignedRead > uint64(len(buf)) {
		alignedRead = uint64(len(buf))
	}

	n, err := unix.Pread(r.fd, buf[:alignedRead], int64(r.pos))
	if err != nil && err != io.EOF {
		return nil, 0, fmt.Errorf("pread at offset %d: %w", r.pos, err)
	}

	actual := int(toRead)
	if n < actual {
		actual = n
	}
	if actual <= 0 {
		return nil, 0, io.EOF
	}

	r.pos += uint64(actual)
	return buf[:actual], actual, nil
}

// Close closes the underlying file descriptor.
func (r *Reader) Close() error {
	return unix.Close(r.fd)
}

// deviceSize determines the size of a block device or regular file.
func deviceSize(fd int, path string) (uint64, error) {
	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		return 0, fmt.Errorf("fstat %s: %w", path, err)
	}

	// Regular file — use stat size.
	if stat.Mode&unix.S_IFMT == unix.S_IFREG {
		return uint64(stat.Size), nil
	}

	// Block device — use BLKGETSIZE64 ioctl.
	if stat.Mode&unix.S_IFMT == unix.S_IFBLK {
		var size uint64
		_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.BLKGETSIZE64, uintptr(unsafe.Pointer(&size)))
		if errno != 0 {
			return 0, fmt.Errorf("ioctl BLKGETSIZE64 on %s: %w", path, errno)
		}
		return size, nil
	}

	// Try seeking to end as a fallback.
	end, err := unix.Seek(fd, 0, unix.SEEK_END)
	if err != nil {
		return 0, fmt.Errorf("cannot determine size of %s: %w", path, err)
	}
	if _, err := unix.Seek(fd, 0, unix.SEEK_SET); err != nil {
		return 0, fmt.Errorf("seek back to start of %s: %w", path, err)
	}
	return uint64(end), nil
}

// alignedAlloc allocates a byte slice whose underlying array is aligned to
// the required boundary for O_DIRECT.
func alignedAlloc(size int) []byte {
	// Allocate extra bytes so we can find an aligned offset within.
	raw := make([]byte, size+alignment)
	addr := uintptr(unsafe.Pointer(&raw[0]))
	offset := (alignment - int(addr%uintptr(alignment))) % alignment
	return raw[offset : offset+size]
}

// FileInfo wraps the minimal information we need.
type FileInfo struct {
	Path string
	Size uint64
}

// Stat returns information about the path without opening it with O_DIRECT.
func Stat(path string) (*FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Path: path,
		Size: uint64(fi.Size()),
	}, nil
}
