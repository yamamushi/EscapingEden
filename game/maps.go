package game

func (gm *GameManager) GetMapChunkByID(id string) *MapChunk {
	for _, chunk := range gm.ChunkCache {
		if chunk.ID == id {
			return chunk
		}
	}
	// Optionally, scan disk for the filename if needed
	return nil
}
