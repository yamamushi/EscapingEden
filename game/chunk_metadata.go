package game

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ChunkMetadata stores lightweight information about map chunks
type ChunkMetadata struct {
	ID           string    `json:"id"`
	GlobalX      int       `json:"global_x"`
	GlobalY      int       `json:"global_y"`
	GlobalZ      int       `json:"global_z"`
	Checksum     string    `json:"checksum"`
	FileSize     int64     `json:"file_size"`
	LastModified time.Time `json:"last_modified"`
	Version      int       `json:"version"`
	IsGenerated  bool      `json:"is_generated"`
}

// ChunkRegistry manages chunk metadata for fast validation
type ChunkRegistry struct {
	Chunks   map[string]*ChunkMetadata `json:"chunks"`
	WorldDir string                    `json:"world_dir"`
	Version  int                       `json:"version"`
}

// NewChunkRegistry creates a new chunk registry
func NewChunkRegistry(worldDir string) *ChunkRegistry {
	return &ChunkRegistry{
		Chunks:   make(map[string]*ChunkMetadata),
		WorldDir: worldDir,
		Version:  1,
	}
}

// LoadRegistry loads the chunk registry from disk
func (cr *ChunkRegistry) LoadRegistry() error {
	registryPath := filepath.Join(cr.WorldDir, "chunk_registry.json")

	data, err := os.ReadFile(registryPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Registry doesn't exist, start fresh
			return nil
		}
		return fmt.Errorf("failed to read registry: %w", err)
	}

	return json.Unmarshal(data, cr)
}

// SaveRegistry saves the chunk registry to disk
func (cr *ChunkRegistry) SaveRegistry() error {
	registryPath := filepath.Join(cr.WorldDir, "chunk_registry.json")

	data, err := json.MarshalIndent(cr, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	return os.WriteFile(registryPath, data, 0644)
}

// ValidateChunk quickly validates a chunk without loading it
func (cr *ChunkRegistry) ValidateChunk(chunkID string) (bool, error) {
	metadata, exists := cr.Chunks[chunkID]
	if !exists {
		return false, nil
	}

	filename := fmt.Sprintf("%s.map", chunkID)
	filePath := filepath.Join(cr.WorldDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Quick validation: check file size and modification time
	if fileInfo.Size() != metadata.FileSize ||
		!fileInfo.ModTime().Equal(metadata.LastModified) {
		// File has changed, need to update metadata
		return false, nil
	}

	return true, nil
}

// UpdateChunkMetadata updates metadata for a chunk
func (cr *ChunkRegistry) UpdateChunkMetadata(chunkID string, globalX, globalY, globalZ int) error {
	filename := fmt.Sprintf("%s.map", chunkID)
	filePath := filepath.Join(cr.WorldDir, filename)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat chunk file: %w", err)
	}

	// Calculate checksum
	checksum, err := cr.calculateChecksum(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	metadata := &ChunkMetadata{
		ID:           chunkID,
		GlobalX:      globalX,
		GlobalY:      globalY,
		GlobalZ:      globalZ,
		Checksum:     checksum,
		FileSize:     fileInfo.Size(),
		LastModified: fileInfo.ModTime(),
		Version:      1,
		IsGenerated:  true,
	}

	cr.Chunks[chunkID] = metadata
	return nil
}

// calculateChecksum calculates MD5 checksum of a file
func (cr *ChunkRegistry) calculateChecksum(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash), nil
}

// GetMissingChunks returns a list of chunks that need to be generated
func (cr *ChunkRegistry) GetMissingChunks(worldX, worldY, worldZ int) []string {
	var missing []string

	for x := 0; x < worldX; x++ {
		for y := 0; y < worldY; y++ {
			for z := 0; z < worldZ; z++ {
				chunkID := fmt.Sprintf("%d-%d-%d", x, y, z)

				valid, err := cr.ValidateChunk(chunkID)
				if err != nil || !valid {
					missing = append(missing, chunkID)
				}
			}
		}
	}

	return missing
}

// GetChunkCount returns the total number of registered chunks
func (cr *ChunkRegistry) GetChunkCount() int {
	return len(cr.Chunks)
}

// GetValidChunkCount returns the number of valid chunks
func (cr *ChunkRegistry) GetValidChunkCount(worldX, worldY, worldZ int) int {
	validCount := 0

	for x := 0; x < worldX; x++ {
		for y := 0; y < worldY; y++ {
			for z := 0; z < worldZ; z++ {
				chunkID := fmt.Sprintf("%d-%d-%d", x, y, z)

				if valid, _ := cr.ValidateChunk(chunkID); valid {
					validCount++
				}
			}
		}
	}

	return validCount
}

// ScanExistingChunks scans the world directory for existing chunk files and updates the registry
func (cr *ChunkRegistry) ScanExistingChunks(worldX, worldY, worldZ int) error {
	scannedCount := 0
	updatedCount := 0

	for x := 0; x < worldX; x++ {
		for y := 0; y < worldY; y++ {
			for z := 0; z < worldZ; z++ {
				chunkID := fmt.Sprintf("%d-%d-%d", x, y, z)
				filename := fmt.Sprintf("%s.map", chunkID)
				filePath := filepath.Join(cr.WorldDir, filename)

				scannedCount++

				// Check if file exists on disk
				if _, err := os.Stat(filePath); err != nil {
					if os.IsNotExist(err) {
						// File doesn't exist, remove from registry if present
						delete(cr.Chunks, chunkID)
						continue
					}
					// Other error, skip this chunk
					continue
				}

				// File exists, check if registry needs updating
				metadata, exists := cr.Chunks[chunkID]
				if !exists {
					// File exists but not in registry, add it
					err := cr.UpdateChunkMetadata(chunkID, x, y, z)
					if err == nil {
						updatedCount++
					}
				} else {
					// File exists and in registry, validate metadata
					fileInfo, err := os.Stat(filePath)
					if err == nil {
						if fileInfo.Size() != metadata.FileSize ||
							!fileInfo.ModTime().Equal(metadata.LastModified) {
							// Metadata is outdated, update it
							err := cr.UpdateChunkMetadata(chunkID, x, y, z)
							if err == nil {
								updatedCount++
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// GetRegistryStats returns statistics about the chunk registry
func (cr *ChunkRegistry) GetRegistryStats(worldX, worldY, worldZ int) (int, int, int) {
	totalExpected := worldX * worldY * worldZ
	registeredCount := len(cr.Chunks)
	validCount := cr.GetValidChunkCount(worldX, worldY, worldZ)

	return totalExpected, registeredCount, validCount
}
