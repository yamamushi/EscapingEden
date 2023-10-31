package game

func (gm *GameManager) GetMapChunkByID(ID string) *MapChunk {
	for _, chunk := range gm.MapChunks {
		if chunk.ID == ID {
			return &chunk
		}
	}
	return nil
}
