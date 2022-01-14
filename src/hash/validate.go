package hash

func ValidateChunkHash(chunkHash string, chunkData []byte) bool {
	return Chunk(chunkData) == chunkHash
}
