package game

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/yamamushi/EscapingEden/logging"
)

// FastWorldValidation performs quick world validation without loading all chunks
func (gm *GameManager) FastWorldValidation() {
	gm.Log.Println(logging.LogInfo, "Performing fast world validation...")

	worldDimensionsString := gm.Config.WorldGen.Dimensions
	worldX, worldY, worldZ := gm.ParseWorldDimensions(worldDimensionsString)

	// Load chunk registry
	err := gm.ChunkRegistry.LoadRegistry()
	if err != nil {
		gm.Log.Println(logging.LogWarn, "Failed to load chunk registry:", err)
	}

	// Scan existing chunks to ensure registry is up-to-date
	gm.Log.Println(logging.LogInfo, "Scanning existing chunks...")
	err = gm.ChunkRegistry.ScanExistingChunks(worldX, worldY, worldZ)
	if err != nil {
		gm.Log.Println(logging.LogWarn, "Failed to scan existing chunks:", err)
	}

	// Get registry statistics
	totalExpected, registeredCount, validCount := gm.ChunkRegistry.GetRegistryStats(worldX, worldY, worldZ)
	gm.Log.Println(logging.LogInfo, fmt.Sprintf("Chunk registry stats: %d expected, %d registered, %d valid",
		totalExpected, registeredCount, validCount))

	// Quick validation using updated registry
	missingChunks := gm.ChunkRegistry.GetMissingChunks(worldX, worldY, worldZ)
	existingChunks := totalExpected - len(missingChunks)

	gm.Log.Println(logging.LogInfo, fmt.Sprintf("World validation: %d/%d chunks exist, %d need generation",
		existingChunks, totalExpected, len(missingChunks)))

	if len(missingChunks) == totalExpected {
		// No world exists, generate it
		gm.Log.Println(logging.LogInfo, "No world found, generating new world...")
		gm.GenerateWorldAsync(worldX, worldY, worldZ)
	} else if len(missingChunks) > 0 {
		// Some chunks missing, regenerate them
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Regenerating %d missing chunks (skipping %d existing)...",
			len(missingChunks), existingChunks))
		gm.RegenerateMissingChunks(missingChunks)
	} else {
		// World is complete, just validate one chunk to ensure loading works
		gm.Log.Println(logging.LogInfo, "World validation complete, all chunks exist. Testing chunk loading...")
		gm.TestChunkLoading()
	}

	gm.WorldValidated = true
	gm.Log.Println(logging.LogInfo, "Fast world validation completed!")
}

// GenerateWorldAsync generates the world asynchronously, skipping existing chunks
func (gm *GameManager) GenerateWorldAsync(worldX, worldY, worldZ int) {
	go func() {
		gm.Log.Println(logging.LogInfo, "Starting async world generation...")

		totalChunks := worldX * worldY * worldZ
		generated := 0
		skipped := 0

		for i := 0; i < worldX; i++ {
			for j := 0; j < worldY; j++ {
				for k := 0; k < worldZ; k++ {
					filename := fmt.Sprintf("%d-%d-%d.map", i, j, k)
					chunkID := fmt.Sprintf("%d-%d-%d", i, j, k)

					// Check if chunk already exists
					if gm.ChunkExists(filename) {
						skipped++
						// Still update registry for existing chunks
						gm.ChunkRegistry.UpdateChunkMetadata(chunkID, i, j, k)
						continue
					}

					// Generate new chunk
					mapChunkID := uuid.New().String()
					mapChunk := gm.CreateMapChunk(255, 255, 255, i, j, k, mapChunkID)

					err := gm.SaveMapChunk(mapChunk, filename, false)
					if err != nil {
						gm.Log.Println(logging.LogError, "Failed to save map chunk:", err.Error())
						continue
					}

					// Update registry
					gm.ChunkRegistry.UpdateChunkMetadata(chunkID, i, j, k)

					generated++
					if (generated+skipped)%100 == 0 {
						gm.Log.Println(logging.LogInfo, fmt.Sprintf("Processed %d/%d chunks (Generated: %d, Skipped: %d)",
							generated+skipped, totalChunks, generated, skipped))
					}
				}
			}
		}

		// Save registry
		gm.ChunkRegistry.SaveRegistry()
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("World generation complete! Generated %d new chunks, skipped %d existing chunks",
			generated, skipped))
	}()
}

// RegenerateMissingChunks regenerates only the missing chunks
func (gm *GameManager) RegenerateMissingChunks(missingChunks []string) {
	go func() {
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Regenerating %d missing chunks...", len(missingChunks)))

		generated := 0
		skipped := 0

		for i, chunkID := range missingChunks {
			// Parse chunk coordinates
			var x, y, z int
			fmt.Sscanf(chunkID, "%d-%d-%d", &x, &y, &z)

			filename := fmt.Sprintf("%s.map", chunkID)

			// Double-check if chunk exists (might have been created since registry was built)
			if gm.ChunkExists(filename) {
				skipped++
				// Update registry for existing chunk
				gm.ChunkRegistry.UpdateChunkMetadata(chunkID, x, y, z)
				continue
			}

			mapChunkID := uuid.New().String()
			mapChunk := gm.CreateMapChunk(255, 255, 255, x, y, z, mapChunkID)

			err := gm.SaveMapChunk(mapChunk, filename, true)
			if err != nil {
				gm.Log.Println(logging.LogError, "Failed to save missing chunk:", err.Error())
				continue
			}

			// Update registry
			gm.ChunkRegistry.UpdateChunkMetadata(chunkID, x, y, z)
			generated++

			if (i+1)%10 == 0 {
				gm.Log.Println(logging.LogInfo, fmt.Sprintf("Processed %d/%d missing chunks (Generated: %d, Skipped: %d)",
					i+1, len(missingChunks), generated, skipped))
			}
		}

		// Save registry
		gm.ChunkRegistry.SaveRegistry()
		gm.Log.Println(logging.LogInfo, fmt.Sprintf("Missing chunk regeneration complete! Generated %d new chunks, skipped %d existing chunks",
			generated, skipped))
	}()
}

// ChunkExists checks if a chunk file already exists on disk
func (gm *GameManager) ChunkExists(filename string) bool {
	fullPath := "./assets/world/" + filename
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// TestChunkLoading tests that chunk loading works without loading all chunks
func (gm *GameManager) TestChunkLoading() {
	gm.Log.Println(logging.LogInfo, "Testing chunk loading...")

	// Test loading just one chunk to verify the system works
	chunk, err := gm.LazyLoader.GetChunkAsync(0, 0, 0)
	if err != nil {
		gm.Log.Println(logging.LogError, "Chunk loading test failed:", err)
		return
	}

	if len(chunk.TileMap) == 0 || len(chunk.TileMap[0]) == 0 || len(chunk.TileMap[0][0]) == 0 {
		gm.Log.Println(logging.LogError, "Loaded chunk is empty!")
		return
	}

	gm.Log.Println(logging.LogInfo, "Chunk loading test successful!")
}

// LoadWorld - kept for backward compatibility but now uses fast validation
func (gm *GameManager) LoadWorld() {
	gm.FastWorldValidation()
}

func (gm *GameManager) ValidateWorldFiles(x, y, z int) error {
	worldDir := "assets/world/"
	// check if world directory exists
	// if not, create it
	_, err := os.Stat(worldDir)
	if err != nil {
		if os.IsNotExist(err) {
			// create directory
			err := os.Mkdir(worldDir, 0755)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	var missingFiles []string
	// Check for all map chunks in range to validate they exist
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < z; k++ {
				mapChunkFilename := strconv.Itoa(i) + "-" + strconv.Itoa(j) + "-" + strconv.Itoa(k) + ".map"
				_, err := os.Stat(worldDir + mapChunkFilename)
				if err != nil {
					if os.IsNotExist(err) {
						missingFiles = append(missingFiles, mapChunkFilename)
					} else {
						return err
					}
				}
			}
		}
	}

	if len(missingFiles) > 0 && len(missingFiles) != (x)*(y)*(z) {
		return errors.New("world files are missing, will not regenerate missing chunks. Aborting")
	}
	if len(missingFiles) == (x)*(y)*(z) {
		gm.Log.Println(logging.LogInfo, "No world files found!")
		return errors.New("empty")
	}
	return nil
}

func (gm *GameManager) ParseWorldDimensions(dimensions string) (int, int, int) {
	stringCollection := strings.Split(dimensions, ",")
	if len(stringCollection) != 3 {
		panic("World dimensions are invalid, please check config!")
	}

	x, err := strconv.Atoi(stringCollection[0])
	if err != nil {
		panic(err)
	}
	if x < 1 {
		panic("World dimensions are invalid, expecting x larger than 0!")
	}
	y, err := strconv.Atoi(stringCollection[1])
	if err != nil {
		panic(err)
	}
	if y < 1 {
		panic("World dimensions are invalid, expecting y larger than 0!")
	}
	z, err := strconv.Atoi(stringCollection[2])
	if err != nil {
		panic(err)
	}
	if z < 1 {
		panic("World dimensions are invalid, expecting z larger than 0!")
	}
	return x, y, z
}
