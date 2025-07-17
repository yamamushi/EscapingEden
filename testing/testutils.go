package testing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

// TestLogger implements a logger for testing that captures output
type TestLogger struct {
	buffer bytes.Buffer
	mutex  sync.RWMutex
}

// NewTestLogger creates a new test logger
func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

// GetTypeID returns the logger type ID
func (tl *TestLogger) GetTypeID() logging.LoggerTypeID {
	return logging.LoggerTypeID_Console
}

// Println logs a message to the test buffer
func (tl *TestLogger) Println(level logging.LogLevel, message string, v ...interface{}) {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	if v == nil {
		fmt.Fprintf(&tl.buffer, "[%s] %s: %s\n", timestamp, level.String(), message)
	} else {
		fmt.Fprintf(&tl.buffer, "[%s] %s: %s %v\n", timestamp, level.String(), message, v)
	}
}

// GetOutput returns the captured log output
func (tl *TestLogger) GetOutput() string {
	tl.mutex.RLock()
	defer tl.mutex.RUnlock()
	return tl.buffer.String()
}

// Clear clears the log buffer
func (tl *TestLogger) Clear() {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	tl.buffer.Reset()
}

// MockDatabase implements a mock database for testing
type MockDatabase struct {
	data   map[string]map[string]interface{}
	mutex  sync.RWMutex
	errors map[string]error // Map operation names to errors for testing error conditions
}

// NewMockDatabase creates a new mock database
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		data:   make(map[string]map[string]interface{}),
		errors: make(map[string]error),
	}
}

// GetTypeID returns the database type ID
func (md *MockDatabase) GetTypeID() edendb.DatabaseTypeID {
	return edendb.DatabaseTypeID_Unknown
}

// GetTypeName returns the database type name
func (md *MockDatabase) GetTypeName() string {
	return "Mock"
}

// Init initializes the mock database
func (md *MockDatabase) Init() error {
	if err, exists := md.errors["Init"]; exists {
		return err
	}
	return nil
}

// AddRecord adds a record to the mock database
func (md *MockDatabase) AddRecord(collection string, value interface{}) error {
	if err, exists := md.errors["AddRecord"]; exists {
		return err
	}

	md.mutex.Lock()
	defer md.mutex.Unlock()

	if md.data[collection] == nil {
		md.data[collection] = make(map[string]interface{})
	}

	// Use a simple key generation for testing
	key := fmt.Sprintf("record_%d", len(md.data[collection]))
	md.data[collection][key] = value

	return nil
}

// UpdateRecord updates a record in the mock database
func (md *MockDatabase) UpdateRecord(collection string, value interface{}) error {
	if err, exists := md.errors["UpdateRecord"]; exists {
		return err
	}

	// For simplicity, just add the record
	return md.AddRecord(collection, value)
}

// UpdateField updates a field in a record
func (md *MockDatabase) UpdateField(collection, field string, value, record interface{}) error {
	if err, exists := md.errors["UpdateField"]; exists {
		return err
	}

	// Simplified implementation for testing
	return nil
}

// AddIfNotExists adds a record if it doesn't exist
func (md *MockDatabase) AddIfNotExists(collection string, value interface{}) error {
	if err, exists := md.errors["AddIfNotExists"]; exists {
		return err
	}

	return md.AddRecord(collection, value)
}

// RemoveRecord removes a record from the collection
func (md *MockDatabase) RemoveRecord(collection string, value interface{}) error {
	if err, exists := md.errors["RemoveRecord"]; exists {
		return err
	}

	md.mutex.Lock()
	defer md.mutex.Unlock()

	if md.data[collection] != nil {
		// Simple removal - remove first matching record
		for key := range md.data[collection] {
			delete(md.data[collection], key)
			break
		}
	}

	return nil
}

// RemoveCollection removes an entire collection
func (md *MockDatabase) RemoveCollection(collection string) error {
	if err, exists := md.errors["RemoveCollection"]; exists {
		return err
	}

	md.mutex.Lock()
	defer md.mutex.Unlock()

	delete(md.data, collection)
	return nil
}

// One finds a single record
func (md *MockDatabase) One(collection, field string, value, output interface{}) error {
	if err, exists := md.errors["One"]; exists {
		return err
	}

	md.mutex.RLock()
	defer md.mutex.RUnlock()

	if md.data[collection] == nil {
		return fmt.Errorf("not found")
	}

	// Simple implementation - return first record
	for _, record := range md.data[collection] {
		// In a real implementation, we'd check the field value
		// For testing, just return the first record
		if account, ok := output.(*messages.Account); ok {
			if mockAccount, ok := record.(*messages.Account); ok {
				*account = *mockAccount
				return nil
			}
		}
		break
	}

	return fmt.Errorf("not found")
}

// FindAll finds all records in a collection
func (md *MockDatabase) FindAll(collection string, output []interface{}) error {
	if err, exists := md.errors["FindAll"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// FindAllByField finds all records by field value
func (md *MockDatabase) FindAllByField(collection, field, searchValue string, output []interface{}) error {
	if err, exists := md.errors["FindAllByField"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// DumpDatabase dumps the entire database
func (md *MockDatabase) DumpDatabase(output []interface{}) error {
	if err, exists := md.errors["DumpDatabase"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// DumpCollection dumps a collection
func (md *MockDatabase) DumpCollection(collection string, output []interface{}) error {
	if err, exists := md.errors["DumpCollection"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// DumpToWriter dumps database to writer
func (md *MockDatabase) DumpToWriter(writer io.Writer) error {
	if err, exists := md.errors["DumpToWriter"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// DumpCollectionToWriter dumps collection to writer
func (md *MockDatabase) DumpCollectionToWriter(collection string, writer io.Writer) error {
	if err, exists := md.errors["DumpCollectionToWriter"]; exists {
		return err
	}

	// Simplified implementation
	return nil
}

// SetError sets an error for a specific operation (for testing error conditions)
func (md *MockDatabase) SetError(operation string, err error) {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	md.errors[operation] = err
}

// ClearErrors clears all set errors
func (md *MockDatabase) ClearErrors() {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	md.errors = make(map[string]error)
}

// TestConfig creates a test configuration
func TestConfig() edenconfig.Config {
	return edenconfig.Config{
		Server: edenconfig.ServerConfig{
			Host:            "localhost",
			Port:            "0", // Use random port for testing
			ShutdownTimeout: 1,
		},
		Logger: edenconfig.LoggerConfig{
			Type: "console",
			Path: "",
		},
		DB: edenconfig.DatabaseConfig{
			Type: "mock",
			Path: ":memory:",
		},
		Discord: edenconfig.DiscordConfig{
			BotToken:              "test_token",
			GuildID:               "test_guild",
			AdminIDs:              []string{"test_admin"},
			RegistrationChannelID: "test_channel",
			RegisteredRoleID:      "test_role",
		},
		WorldGen: edenconfig.WorldGenConfig{
			Dimensions: "5,5,5",
			ChunkSize:  32,
		},
	}
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "eden_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	return dir
}

// CreateTestConfigFile creates a test configuration file
func CreateTestConfigFile(t *testing.T, dir string) string {
	configPath := filepath.Join(dir, "test.conf")

	configContent := `[server]
host = "localhost"
port = "0"
shutdown_timeout = 1

[logging]
type = "console"
path = ""

[database]
type = "bolt"
path = ":memory:"

[discord]
bot_token = "test_token"
guild_id = "test_guild"
admin_ids = ["test_admin"]
registration_channel_id = "test_channel"
registered_role_id = "test_role"

[worldgen]
dimensions = "5,5,5"
chunk_size = 32
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	return configPath
}

// MockConnection implements a mock network connection for testing
type MockConnection struct {
	readBuffer  bytes.Buffer
	writeBuffer bytes.Buffer
	closed      bool
	mutex       sync.RWMutex
}

// NewMockConnection creates a new mock connection
func NewMockConnection() *MockConnection {
	return &MockConnection{}
}

// Read reads from the mock connection
func (mc *MockConnection) Read(b []byte) (n int, err error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	if mc.closed {
		return 0, io.EOF
	}

	return mc.readBuffer.Read(b)
}

// Write writes to the mock connection
func (mc *MockConnection) Write(b []byte) (n int, err error) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if mc.closed {
		return 0, io.ErrClosedPipe
	}

	return mc.writeBuffer.Write(b)
}

// Close closes the mock connection
func (mc *MockConnection) Close() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.closed = true
	return nil
}

// LocalAddr returns a mock local address
func (mc *MockConnection) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

// RemoteAddr returns a mock remote address
func (mc *MockConnection) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
}

// SetDeadline sets the deadline (no-op for mock)
func (mc *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the read deadline (no-op for mock)
func (mc *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the write deadline (no-op for mock)
func (mc *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

// WriteToInput writes data to the input buffer (simulates incoming data)
func (mc *MockConnection) WriteToInput(data []byte) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.readBuffer.Write(data)
}

// GetOutput returns the output buffer contents
func (mc *MockConnection) GetOutput() string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	return mc.writeBuffer.String()
}

// ClearOutput clears the output buffer
func (mc *MockConnection) ClearOutput() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.writeBuffer.Reset()
}

// TestHelper provides common testing utilities
type TestHelper struct {
	t       *testing.T
	logger  *TestLogger
	db      *MockDatabase
	config  edenconfig.Config
	tempDir string
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir := CreateTempDir(t)

	return &TestHelper{
		t:       t,
		logger:  NewTestLogger(),
		db:      NewMockDatabase(),
		config:  TestConfig(),
		tempDir: tempDir,
	}
}

// AssertNoError asserts that an error is nil
func (th *TestHelper) AssertNoError(err error) {
	if err != nil {
		th.t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError asserts that an error is not nil
func (th *TestHelper) AssertError(err error) {
	if err == nil {
		th.t.Fatal("Expected an error, got nil")
	}
}

// AssertEqual asserts that two values are equal
func (th *TestHelper) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		th.t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertContains asserts that a string contains a substring
func (th *TestHelper) AssertContains(str, substr string) {
	if !strings.Contains(str, substr) {
		th.t.Fatalf("Expected string to contain %q, got: %q", substr, str)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func (th *TestHelper) AssertNotContains(str, substr string) {
	if strings.Contains(str, substr) {
		th.t.Fatalf("Expected string to not contain %q, got: %q", substr, str)
	}
}

// GetLogger returns the test logger
func (th *TestHelper) GetLogger() *TestLogger {
	return th.logger
}

// GetDB returns the mock database
func (th *TestHelper) GetDB() *MockDatabase {
	return th.db
}

// GetConfig returns the test configuration
func (th *TestHelper) GetConfig() edenconfig.Config {
	return th.config
}

// GetTempDir returns the temporary directory
func (th *TestHelper) GetTempDir() string {
	return th.tempDir
}

// WithTimeout runs a function with a timeout
func (th *TestHelper) WithTimeout(timeout time.Duration, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		fn(ctx)
	}()

	select {
	case <-done:
		// Function completed
	case <-ctx.Done():
		th.t.Fatal("Test timed out")
	}
}
