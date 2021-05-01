package schema

import (
	"errors"
	"io"
	"math"
	"os"
)

var (
	// ErrInvalidFile represent a invalid file
	ErrInvalidFile = errors.New("invalid file")
	// ErrFileNoChunks represent that a file has no any chunks
	ErrFileNoChunks = errors.New("file has no any chunks")
	// ErrInvalidSeekWhence represent invalid seek whence
	ErrInvalidSeekWhence = errors.New("invalid seek whence")
	// ErrNegativePosition represent negative position
	ErrNegativePosition = errors.New("negative read position")
)

type fileReader struct {
	file               *File
	rootPath           *string
	currentChunkReader *os.File
	totalChunkNumber   int
	currentChunkNumber int
	alreadyReadCount   int
}

func NewFileReader(file *File, rootPath *string) (*fileReader, error) {
	if file == nil {
		return nil, ErrInvalidFile
	}

	var (
		err              error
		firstChunk       *Chunk
		chunkReader      *os.File
		totalChunkNumber int
	)

	if totalChunkNumber, err = file.LastChunkNumber(); err != nil {
		return nil, err
	}

	if totalChunkNumber == 0 {
		return nil, ErrFileNoChunks
	}

	if firstChunk, err = file.ChunkByNumber(1); err != nil {
		return nil, err
	}

	if chunkReader, err = firstChunk.Reader(rootPath); err != nil {
		return nil, err
	}

	return &fileReader{
		file:               file,
		currentChunkReader: chunkReader,
		rootPath:           rootPath,
		currentChunkNumber: 1,
		totalChunkNumber:   totalChunkNumber,
	}, nil
}

func (fr *fileReader) Read(p []byte) (readCount int, err error) {
	if fr.alreadyReadCount >= int(fr.file.Size) {
		_ = fr.currentChunkReader.Close()
		return 0, io.EOF
	}
	defer func() { fr.alreadyReadCount += readCount }()
	readCount, err = fr.currentChunkReader.Read(p)
	if err != nil && err == io.EOF {
		_ = fr.currentChunkReader.Close()
		fr.currentChunkNumber++
		var nextChunk *Chunk
		if nextChunk, err = fr.file.ChunkByNumber(fr.currentChunkNumber); err != nil {
			return
		}
		if fr.currentChunkReader, err = nextChunk.Reader(fr.rootPath); err != nil {
			return readCount, err
		}
		return readCount, nil
	}
	return readCount, err
}

func (fr *fileReader) Seek(offset int64, whence int) (abs int64, err error) {
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(fr.alreadyReadCount) + offset
	case io.SeekEnd:
		abs = fr.file.Size + offset
	default:
		return 0, ErrInvalidSeekWhence
	}
	if abs < 0 {
		return 0, ErrNegativePosition
	}
	if abs >= fr.file.Size {
		fr.alreadyReadCount = int(abs)
		fr.currentChunkNumber = fr.totalChunkNumber
		return abs, nil
	}
	var (
		currentChunk       *Chunk
		currentChunkReader *os.File
		currentChunkNumber = int(math.Ceil(float64(abs) / float64(fr.file.ChunkSize)))
	)

	if abs%fr.file.ChunkSize == 0 {
		currentChunkNumber++
	}

	if currentChunkNumber == fr.currentChunkNumber {
		currentChunkReader = fr.currentChunkReader
	} else {
		if currentChunk, err = fr.file.ChunkByNumber(currentChunkNumber); err != nil {
			return 0, nil
		}
		if currentChunkReader, err = currentChunk.Reader(fr.rootPath); err != nil {
			return 0, err
		}
	}
	if _, err = currentChunkReader.Seek(abs%fr.file.ChunkSize, io.SeekStart); err != nil {
		return 0, err
	}
	if currentChunkNumber != fr.currentChunkNumber {
		_ = fr.currentChunkReader.Close()
	}
	fr.currentChunkReader = currentChunkReader
	fr.currentChunkNumber = currentChunkNumber
	fr.alreadyReadCount = int(abs)
	return abs, nil
}
