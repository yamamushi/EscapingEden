package game

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// ConfigManager manages configuration files and hot-reload functionality
type ConfigManager struct {
	actionRegistry *ActionRegistry
	watchers       map[string]*FileWatcher
	mutex          sync.RWMutex
	log            logging.LoggerType
	enabled        bool
}

// FileWatcher watches a single file for changes
type FileWatcher struct {
	filePath    string
	lastModTime time.Time
	callback    func(string) error
	stopChan    chan struct{}
	running     bool
	mutex       sync.Mutex
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(actionRegistry *ActionRegistry, log logging.LoggerType) *ConfigManager {
	return &ConfigManager{
		actionRegistry: actionRegistry,
		watchers:       make(map[string]*FileWatcher),
		log:            log,
		enabled:        true,
	}
}

// StartWatching starts watching configuration files for changes
func (cm *ConfigManager) StartWatching() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if !cm.enabled {
		return fmt.Errorf("config manager is disabled")
	}

	// Watch actions.json file
	actionsConfigPath := "config/actions.json"
	if _, err := os.Stat(actionsConfigPath); err == nil {
		err := cm.watchFile(actionsConfigPath, cm.handleActionsConfigChange)
		if err != nil {
			return fmt.Errorf("failed to watch actions config: %v", err)
		}
		cm.log.Println(logging.LogInfo, fmt.Sprintf("Started watching %s for changes", actionsConfigPath))
	}

	return nil
}

// StopWatching stops all file watchers
func (cm *ConfigManager) StopWatching() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for filePath, watcher := range cm.watchers {
		watcher.Stop()
		delete(cm.watchers, filePath)
		cm.log.Println(logging.LogInfo, fmt.Sprintf("Stopped watching %s", filePath))
	}
}

// watchFile starts watching a specific file for changes
func (cm *ConfigManager) watchFile(filePath string, callback func(string) error) error {
	// Get initial file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", filePath, err)
	}

	// Create file watcher
	watcher := &FileWatcher{
		filePath:    filePath,
		lastModTime: fileInfo.ModTime(),
		callback:    callback,
		stopChan:    make(chan struct{}),
	}

	// Start watching in a goroutine
	go watcher.watch(cm.log)

	// Store watcher
	cm.watchers[filePath] = watcher

	return nil
}

// handleActionsConfigChange handles changes to the actions.json file
func (cm *ConfigManager) handleActionsConfigChange(filePath string) error {
	cm.log.Println(logging.LogInfo, fmt.Sprintf("Actions configuration file changed: %s", filePath))

	// Create backup of current configuration
	backup := cm.createActionRegistryBackup()

	// Attempt to reload configuration
	err := cm.actionRegistry.ReloadFromFile()
	if err != nil {
		cm.log.Println(logging.LogError, fmt.Sprintf("Failed to reload actions configuration: %v", err))

		// Restore backup
		cm.restoreActionRegistryBackup(backup)
		cm.log.Println(logging.LogWarn, "Restored previous actions configuration due to reload failure")

		return fmt.Errorf("configuration reload failed: %v", err)
	}

	cm.log.Println(logging.LogInfo, "Actions configuration reloaded successfully")
	return nil
}

// createActionRegistryBackup creates a backup of the current action registry
func (cm *ConfigManager) createActionRegistryBackup() ActionRegistryBackup {
	return ActionRegistryBackup{
		TickRate:     cm.actionRegistry.GetTickRate(),
		MaxQueueSize: cm.actionRegistry.GetMaxQueueSize(),
		Actions:      cm.copyActionDefinitions(),
	}
}

// restoreActionRegistryBackup restores the action registry from a backup
func (cm *ConfigManager) restoreActionRegistryBackup(backup ActionRegistryBackup) {
	cm.actionRegistry.mutex.Lock()
	defer cm.actionRegistry.mutex.Unlock()

	cm.actionRegistry.TickRate = backup.TickRate
	cm.actionRegistry.MaxQueueSize = backup.MaxQueueSize
	cm.actionRegistry.Actions = backup.Actions
}

// copyActionDefinitions creates a deep copy of action definitions
func (cm *ConfigManager) copyActionDefinitions() map[string]*ActionDefinition {
	cm.actionRegistry.mutex.RLock()
	defer cm.actionRegistry.mutex.RUnlock()

	copy := make(map[string]*ActionDefinition)
	for actionType, definition := range cm.actionRegistry.Actions {
		copy[actionType] = &ActionDefinition{
			TickCost:      definition.TickCost,
			Description:   definition.Description,
			CanQueue:      definition.CanQueue,
			Interruptible: definition.Interruptible,
			Category:      definition.Category,
		}
	}
	return copy
}

// ValidateConfiguration validates a configuration file without loading it
func (cm *ConfigManager) ValidateConfiguration(filePath string) error {
	// Create a temporary registry for validation
	tempRegistry := NewActionRegistry(cm.log)
	err := tempRegistry.LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	cm.log.Println(logging.LogInfo, fmt.Sprintf("Configuration file %s is valid", filePath))
	return nil
}

// GetWatchedFiles returns a list of currently watched files
func (cm *ConfigManager) GetWatchedFiles() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	files := make([]string, 0, len(cm.watchers))
	for filePath := range cm.watchers {
		files = append(files, filePath)
	}
	return files
}

// SetEnabled enables or disables the configuration manager
func (cm *ConfigManager) SetEnabled(enabled bool) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.enabled = enabled
	cm.log.Println(logging.LogInfo, fmt.Sprintf("ConfigManager enabled: %t", enabled))

	if !enabled {
		// Stop all watchers when disabled
		for _, watcher := range cm.watchers {
			watcher.Stop()
		}
	}
}

// IsEnabled returns whether the configuration manager is enabled
func (cm *ConfigManager) IsEnabled() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.enabled
}

// ActionRegistryBackup holds a backup of action registry configuration
type ActionRegistryBackup struct {
	TickRate     int
	MaxQueueSize int
	Actions      map[string]*ActionDefinition
}

// watch monitors a file for changes (runs in its own goroutine)
func (fw *FileWatcher) watch(log logging.LoggerType) {
	fw.mutex.Lock()
	fw.running = true
	fw.mutex.Unlock()

	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fw.checkForChanges(log)
		case <-fw.stopChan:
			fw.mutex.Lock()
			fw.running = false
			fw.mutex.Unlock()
			return
		}
	}
}

// checkForChanges checks if the file has been modified
func (fw *FileWatcher) checkForChanges(log logging.LoggerType) {
	fileInfo, err := os.Stat(fw.filePath)
	if err != nil {
		log.Println(logging.LogWarn, fmt.Sprintf("Failed to check file %s: %v", fw.filePath, err))
		return
	}

	fw.mutex.Lock()
	lastModTime := fw.lastModTime
	fw.mutex.Unlock()

	if fileInfo.ModTime().After(lastModTime) {
		fw.mutex.Lock()
		fw.lastModTime = fileInfo.ModTime()
		fw.mutex.Unlock()

		// File has been modified, call callback
		if fw.callback != nil {
			err := fw.callback(fw.filePath)
			if err != nil {
				log.Println(logging.LogError, fmt.Sprintf("File change callback failed for %s: %v", fw.filePath, err))
			}
		}
	}
}

// Stop stops the file watcher
func (fw *FileWatcher) Stop() {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	if fw.running {
		close(fw.stopChan)
		fw.running = false
	}
}

// IsRunning returns whether the file watcher is running
func (fw *FileWatcher) IsRunning() bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	return fw.running
}

// GetLastModTime returns the last modification time of the watched file
func (fw *FileWatcher) GetLastModTime() time.Time {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	return fw.lastModTime
}
