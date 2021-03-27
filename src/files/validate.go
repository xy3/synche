package files

func ValidateChunkHash(chunkHash string, chunkData []byte) bool {
	return HashChunk(chunkData) == chunkHash
}