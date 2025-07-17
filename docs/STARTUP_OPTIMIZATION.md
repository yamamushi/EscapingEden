# Server Startup Optimization - Skip Existing Map Chunks

This document describes the optimizations implemented to prevent the server from attempting to generate map chunks that already exist on disk, eliminating the "failed to save map" errors during startup.

## Problem Description

Previously, during server startup, the world generation system would:
1. Attempt to generate ALL chunks in the world dimensions
2. Try to save each chunk to disk
3. Fail with "file already exists" errors for chunks that were already present
4. Continue processing, generating unnecessary chunks and logging errors

This resulted in:
- Slow startup times due to unnecessary chunk generation
- Confusing error messages about failed saves
- Wasted CPU and memory resources
- Poor user experience during server restarts

## Solution Overview

The optimization implements a **scan-before-generate** approach:

1. **Chunk Registry Scanning**: Scan existing chunk files on disk
2. **Registry Synchronization**: Update the chunk registry with current file state
3. **Selective Generation**: Only generate chunks that don't exist
4. **Improved Logging**: Clear progress reporting with skip counts

## Key Optimizations Implemented

### 1. Chunk Existence Checking

Added `ChunkExists()` function to check if a chunk file already exists:

```go
func (gm *GameManager) ChunkExists(filename string) bool {
    fullPath := "./assets/world/" + filename
    _, err := os.Stat(fullPath)
    return !os.IsNotExist(err)
}
```

### 2. Enhanced World Generation

Modified `GenerateWorldAsync()` to skip existing chunks:

**Before:**
```go
// Always generate and try to save
mapChunk := gm.CreateMapChunk(255, 255, 255, i, j, k, mapChunkID)
err := gm.SaveMapChunk(mapChunk, filename, false)
if err != nil {
    gm.Log.Println(logging.LogError, "Failed to save map chunk:", err.Error())
    continue
}
```

**After:**
```go
// Check if chunk already exists
if gm.ChunkExists(filename) {
    skipped++
    // Still update registry for existing chunks
    gm.ChunkRegistry.UpdateChunkMetadata(chunkID, i, j, k)
    continue
}

// Only generate if chunk doesn't exist
mapChunk := gm.CreateMapChunk(255, 255, 255, i, j, k, mapChunkID)
err := gm.SaveMapChunk(mapChunk, filename, false)
```

### 3. Registry Synchronization

Added `ScanExistingChunks()` to synchronize the registry with disk state:

```go
func (cr *ChunkRegistry) ScanExistingChunks(worldX, worldY, worldZ int) error {
    // Scan all expected chunk locations
    // Update registry for existing files
    // Remove registry entries for missing files
    // Update metadata for changed files
}
```

### 4. Improved Missing Chunk Regeneration

Enhanced `RegenerateMissingChunks()` with double-checking:

```go
// Double-check if chunk exists (might have been created since registry was built)
if gm.ChunkExists(filename) {
    skipped++
    // Update registry for existing chunk
    gm.ChunkRegistry.UpdateChunkMetadata(chunkID, x, y, z)
    continue
}
```

### 5. Enhanced Progress Reporting

Improved logging to show both generated and skipped chunks:

```go
gm.Log.Println(logging.LogInfo, fmt.Sprintf("Processed %d/%d chunks (Generated: %d, Skipped: %d)", 
    generated+skipped, totalChunks, generated, skipped))
```

## Performance Benefits

### Startup Time Reduction
- **Before**: Generate all chunks regardless of existence
- **After**: Only generate missing chunks
- **Improvement**: Proportional to existing chunk percentage

### Resource Usage
- **CPU**: Dramatically reduced for existing worlds
- **Memory**: Lower peak usage during generation
- **Disk I/O**: Reduced write operations

### Error Elimination
- **Before**: Multiple "failed to save map" errors
- **After**: Clean startup with informative progress messages

## Example Startup Logs

### Before Optimization
```
[INFO] Starting async world generation...
[ERROR] Failed to save map chunk: file 0-0-0.map already exists
[ERROR] Failed to save map chunk: file 0-0-1.map already exists
[ERROR] Failed to save map chunk: file 0-1-0.map already exists
... (hundreds of errors)
[INFO] World generation complete! Generated 0 chunks
```

### After Optimization
```
[INFO] Performing fast world validation...
[INFO] Scanning existing chunks...
[INFO] Chunk registry stats: 1000 expected, 800 registered, 750 valid
[INFO] World validation: 750/1000 chunks exist, 250 need generation
[INFO] Regenerating 250 missing chunks (skipping 750 existing)...
[INFO] Processed 100/250 missing chunks (Generated: 95, Skipped: 5)
[INFO] Missing chunk regeneration complete! Generated 245 new chunks, skipped 5 existing chunks
```

## Technical Implementation Details

### Registry-Based Validation

The system now uses a two-phase approach:

1. **Registry Loading**: Load existing chunk metadata
2. **Disk Scanning**: Verify registry against actual files
3. **Synchronization**: Update registry with current state
4. **Generation Planning**: Determine what needs to be created

### Chunk Metadata Tracking

Enhanced chunk registry with:
- File size validation
- Modification time tracking
- Checksum verification
- Generation status flags

### Error Handling

Improved error handling for:
- Missing registry files (start fresh)
- Corrupted chunk files (regenerate)
- Disk access issues (graceful degradation)
- Registry synchronization failures (continue with warnings)

## Configuration Impact

### World Dimensions
The optimization scales with world size:
- **Small worlds** (10x10x3): Minimal impact, already fast
- **Medium worlds** (50x50x3): Significant improvement
- **Large worlds** (100x100x5): Dramatic startup time reduction

### Existing Worlds
- **Fresh installations**: No change in behavior
- **Existing worlds**: Major startup time improvement
- **Partial worlds**: Intelligent gap filling

## Monitoring and Debugging

### New Log Messages
- Chunk scanning progress
- Registry synchronization status
- Generation vs. skip statistics
- Registry validation results

### Registry Statistics
```go
totalExpected, registeredCount, validCount := gm.ChunkRegistry.GetRegistryStats(worldX, worldY, worldZ)
```

### Debug Information
- Chunk existence verification
- Registry consistency checks
- File system validation
- Generation decision logging

## Future Enhancements

### Planned Improvements
1. **Parallel Scanning**: Multi-threaded chunk existence checking
2. **Incremental Updates**: Only scan changed directories
3. **Checksum Validation**: Verify chunk integrity during scan
4. **Background Sync**: Periodic registry updates during runtime

### Extensibility
The optimization framework supports:
- Custom chunk validation rules
- Pluggable storage backends
- Advanced caching strategies
- Distributed world generation

## Troubleshooting

### Common Issues

#### Registry Out of Sync
**Symptoms**: Chunks exist but marked as missing
**Solution**: Registry automatically rescans on startup

#### Partial Generation Failures
**Symptoms**: Some chunks generated, others skipped unexpectedly
**Solution**: Check disk space and file permissions

#### Performance Regression
**Symptoms**: Startup slower than expected
**Solution**: Verify chunk registry isn't corrupted

### Debug Commands
```bash
# Check world directory
ls -la ./assets/world/*.map | wc -l

# Verify registry file
cat ./assets/world/chunk_registry.json | jq '.chunks | length'

# Monitor startup logs
tail -f server.log | grep -E "(Generated|Skipped|chunks)"
```

## Migration Notes

### Existing Installations
- No manual migration required
- Registry automatically created/updated
- Existing chunks preserved
- No data loss risk

### Backward Compatibility
- Old registry formats supported
- Graceful degradation for missing registry
- Compatible with existing world files
- No breaking changes to APIs

This optimization significantly improves server startup performance while maintaining full compatibility with existing worlds and providing better visibility into the world generation process.