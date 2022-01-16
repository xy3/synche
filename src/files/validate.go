package files

func ValidateChunkHash(chunkHash string, chunkData []byte) bool {
	return ChunkHash(chunkData) == chunkHash
}
