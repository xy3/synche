package data

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
)

type Chunk struct {
	Hash  string
	Num   int64
	Bytes *[]byte
}

func NewChunk(num int64, bytes *[]byte) *Chunk {
	return &Chunk{Hash: hash.Chunk(*bytes), Num: num, Bytes: bytes}
}
