package game

func (gm *GameManager) GetMapChunkByID(ID string) *MapChunk {
	for _, chunk := range gm.MapChunks {
		if chunk.ID == ID {
			return &chunk
		}
	}
	return &gm.MapChunks[0] // This is the default map chunk, if the map chunk is not found.
}
