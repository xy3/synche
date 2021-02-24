package data

type Chunk struct {
	Hash  string
	Num   int64
	Bytes *[]byte
}

func NewChunkWithHash(hash string, num int64, bytes *[]byte) *Chunk {
	return &Chunk{Hash: hash, Num: num, Bytes: bytes}
}

func NewChunk(num int64, bytes *[]byte) *Chunk {
	return &Chunk{Hash: DefaultChunkHashFunc(*bytes), Num: num, Bytes: bytes}
}

