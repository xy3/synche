package files

import (
	"github.com/xy3/synche/src/hash"
)

type Chunk struct {
	Hash  string
	Num   int64
	Bytes *[]byte
}

func NewChunk(num int64, bytes *[]byte) *Chunk {
	return &Chunk{Hash: hash.Chunk(*bytes), Num: num, Bytes: bytes}
}
