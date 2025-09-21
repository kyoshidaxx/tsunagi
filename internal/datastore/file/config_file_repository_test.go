package datastore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	c "github.com/kyoshidaxx/tsunagi/internal/domain/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnv loads .env.test file for testing
func setupTestEnv(t *testing.T) {
	// Load .env.test from project root (go up 3 levels from internal/datastore/file/)
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		t.Logf("Warning: Could not load .env.test file: %v", err)
	}
}

func TestNewConfigFileRepository(t *testing.T) {
	// Test with a relative path
	repo := NewConfigFileRepository(".tsunagi/test-config.json")

	// Verify it's a configFileRepository instance
	configRepo, ok := repo.(*configFileRepository)
	require.True(t, ok, "Should return configFileRepository instance")

	// Verify the file path is constructed correctly
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	expectedPath := filepath.Join(homeDir, ".tsunagi/test-config.json")
	assert.Equal(t, expectedPath, configRepo.filePath)
}

func TestNewConfigFileRepository_WithEnvVariable(t *testing.T) {
	// Load test environment variables
	setupTestEnv(t)

	// Create repository using environment variable (like in cmd/add.go)
	repo := NewConfigFileRepository(os.Getenv("CONFIG_FILE_PATH"))

	// Verify it's a configFileRepository instance
	configRepo, ok := repo.(*configFileRepository)
	require.True(t, ok, "Should return configFileRepository instance")

	// Verify the file path is constructed correctly
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	expectedPath := filepath.Join(homeDir, configFilePath)
	assert.Equal(t, expectedPath, configRepo.filePath)
}

func TestConfigFileRepository_Save_NewFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test-config.json")

	// Create repository with temp file path
	repo := &configFileRepository{filePath: testFilePath}

	// Test data
	config := c.ConfigParam{
		Name:         "test-config",
		Port:         50000,
		ProjectName:  "test-project",
		Region:       "asia-northeast1",
		InstanceName: "test-instance",
	}

	// Save configuration
	err := repo.Save(config)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, testFilePath)

	// Verify file contents
	data, err := os.ReadFile(testFilePath)
	require.NoError(t, err)

	var configs []c.ConfigParam
	err = json.Unmarshal(data, &configs)
	require.NoError(t, err)

	require.Len(t, configs, 1)
	assert.Equal(t, config, configs[0])
}

func TestConfigFileRepository_Save_ExistingFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test-config.json")

	// Create repository with temp file path
	repo := &configFileRepository{filePath: testFilePath}

	// Create initial configuration
	initialConfig := c.ConfigParam{
		Name:         "initial-config",
		Port:         50001,
		ProjectName:  "initial-project",
		Region:       "asia-northeast1",
		InstanceName: "initial-instance",
	}

	// Save initial configuration
	err := repo.Save(initialConfig)
	require.NoError(t, err)

	// Add another configuration
	additionalConfig := c.ConfigParam{
		Name:         "additional-config",
		Port:         50002,
		ProjectName:  "additional-project",
		Region:       "us-central1",
		InstanceName: "additional-instance",
	}

	// Save additional configuration
	err = repo.Save(additionalConfig)
	require.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(testFilePath)
	require.NoError(t, err)

	var configs []c.ConfigParam
	err = json.Unmarshal(data, &configs)
	require.NoError(t, err)

	require.Len(t, configs, 2)
	assert.Equal(t, initialConfig, configs[0])
	assert.Equal(t, additionalConfig, configs[1])
}

func TestConfigFileRepository_Save_EmptyFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "empty-config.json")

	// Create an empty file
	err := os.WriteFile(testFilePath, []byte{}, 0644)
	require.NoError(t, err)

	// Create repository with temp file path
	repo := &configFileRepository{filePath: testFilePath}

	// Test data
	config := c.ConfigParam{
		Name:         "test-config",
		Port:         50000,
		ProjectName:  "test-project",
		Region:       "asia-northeast1",
		InstanceName: "test-instance",
	}

	// Save configuration
	err = repo.Save(config)
	require.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(testFilePath)
	require.NoError(t, err)

	var configs []c.ConfigParam
	err = json.Unmarshal(data, &configs)
	require.NoError(t, err)

	require.Len(t, configs, 1)
	assert.Equal(t, config, configs[0])
}

func TestConfigFileRepository_Save_CreateDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "nested", "directory")
	testFilePath := filepath.Join(testDir, "test-config.json")

	// Create repository with nested directory path
	repo := &configFileRepository{filePath: testFilePath}

	// Test data
	config := c.ConfigParam{
		Name:         "test-config",
		Port:         50000,
		ProjectName:  "test-project",
		Region:       "asia-northeast1",
		InstanceName: "test-instance",
	}

	// Save configuration (should create directory)
	err := repo.Save(config)
	require.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, testDir)

	// Verify file was created
	assert.FileExists(t, testFilePath)
}

func TestConfigFileRepository_checkConfigFileExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.json")
	nonExistingFile := filepath.Join(tempDir, "non-existing.json")

	// Create a file
	err := os.WriteFile(existingFile, []byte("{}"), 0644)
	require.NoError(t, err)

	// Test with existing file
	repo1 := &configFileRepository{filePath: existingFile}
	assert.True(t, repo1.checkConfigFileExists())

	// Test with non-existing file
	repo2 := &configFileRepository{filePath: nonExistingFile}
	assert.False(t, repo2.checkConfigFileExists())
}

func TestConfigFileRepository_loadConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test-config.json")

	// Create repository
	repo := &configFileRepository{filePath: testFilePath}

	// Test loading non-existing file
	_, err := repo.loadConfigFile()
	assert.Error(t, err)

	// Create test data
	expectedConfigs := []c.ConfigParam{
		{
			Name:         "config1",
			Port:         50001,
			ProjectName:  "project1",
			Region:       "asia-northeast1",
			InstanceName: "instance1",
		},
		{
			Name:         "config2",
			Port:         50002,
			ProjectName:  "project2",
			Region:       "us-central1",
			InstanceName: "instance2",
		},
	}

	// Write test data to file
	data, err := json.MarshalIndent(expectedConfigs, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(testFilePath, data, 0644)
	require.NoError(t, err)

	// Test loading existing file
	configs, err := repo.loadConfigFile()
	require.NoError(t, err)
	assert.Equal(t, expectedConfigs, configs)
}

func TestConfigFileRepository_loadConfigFile_EmptyFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "empty-config.json")

	// Create empty file
	err := os.WriteFile(testFilePath, []byte{}, 0644)
	require.NoError(t, err)

	// Create repository
	repo := &configFileRepository{filePath: testFilePath}

	// Test loading empty file
	configs, err := repo.loadConfigFile()
	require.NoError(t, err)
	assert.Empty(t, configs)
}

func TestConfigFileRepository_loadConfigFile_InvalidJSON(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "invalid-config.json")

	// Create file with invalid JSON
	err := os.WriteFile(testFilePath, []byte("invalid json content"), 0644)
	require.NoError(t, err)

	// Create repository
	repo := &configFileRepository{filePath: testFilePath}

	// Test loading invalid JSON
	_, err = repo.loadConfigFile()
	assert.Error(t, err)
}

func TestConfigFileRepository_createConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "new", "directory")
	testFilePath := filepath.Join(testDir, "new-config.json")

	// Create repository
	repo := &configFileRepository{filePath: testFilePath}

	// Test creating file in new directory
	err := repo.createConfigFile()
	require.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, testDir)

	// Verify file was created
	assert.FileExists(t, testFilePath)
}

func TestConfigFileRepository_createConfigFile_ExistingDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "existing-dir-config.json")

	// Create repository
	repo := &configFileRepository{filePath: testFilePath}

	// Test creating file in existing directory
	err := repo.createConfigFile()
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, testFilePath)
}

func TestConfigFileRepository_Save_JSONMarshalError(t *testing.T) {
	// This test is difficult to implement without modifying the code
	// because json.MarshalIndent rarely fails with valid data
	// In a real scenario, you might want to use dependency injection
	// or interfaces to make this testable

	t.Log("JSON marshal error test would require code modification for full coverage")
}

// Benchmark tests
func BenchmarkConfigFileRepository_Save(b *testing.B) {
	tempDir := b.TempDir()
	testFilePath := filepath.Join(tempDir, "benchmark-config.json")

	repo := &configFileRepository{filePath: testFilePath}

	config := c.ConfigParam{
		Name:         "benchmark-config",
		Port:         50000,
		ProjectName:  "benchmark-project",
		Region:       "asia-northeast1",
		InstanceName: "benchmark-instance",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.Save(config)
		require.NoError(b, err)
	}
}

func BenchmarkConfigFileRepository_loadConfigFile(b *testing.B) {
	tempDir := b.TempDir()
	testFilePath := filepath.Join(tempDir, "benchmark-load-config.json")

	// Create test data
	configs := make([]c.ConfigParam, 100)
	for i := 0; i < 100; i++ {
		configs[i] = c.ConfigParam{
			Name:         "config" + string(rune(i)),
			Port:         50000 + i,
			ProjectName:  "project" + string(rune(i)),
			Region:       "asia-northeast1",
			InstanceName: "instance" + string(rune(i)),
		}
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	require.NoError(b, err)
	err = os.WriteFile(testFilePath, data, 0644)
	require.NoError(b, err)

	repo := &configFileRepository{filePath: testFilePath}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.loadConfigFile()
		require.NoError(b, err)
	}
}

// Integration test that simulates the actual command usage
func TestConfigFileRepository_Integration_WithEnvVariable(t *testing.T) {
	// Load test environment variables
	setupTestEnv(t)

	// Get config file path from environment
	configFilePath := os.Getenv("CONFIG_FILE_PATH")

	// Skip test if no config file path is set (empty environment variable)
	if configFilePath == "" {
		t.Skip("CONFIG_FILE_PATH not set, skipping integration test")
	}

	// Create repository using environment variable (like in cmd/add.go)
	repo := NewConfigFileRepository(configFilePath)

	// Test saving configuration
	config := c.ConfigParam{
		Name:         "integration-test-config",
		Port:         50000,
		ProjectName:  "integration-test-project",
		Region:       "asia-northeast1",
		InstanceName: "integration-test-instance",
	}

	// Save configuration
	err := repo.Save(config)
	require.NoError(t, err)

	// Verify the file was created in the expected location
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	expectedFilePath := filepath.Join(homeDir, configFilePath)
	assert.FileExists(t, expectedFilePath)

	// Verify file contents
	data, err := os.ReadFile(expectedFilePath)
	require.NoError(t, err)

	var configs []c.ConfigParam
	err = json.Unmarshal(data, &configs)
	require.NoError(t, err)

	require.Len(t, configs, 1)
	assert.Equal(t, config, configs[0])

	// Clean up
	os.Remove(expectedFilePath)
}

func TestNewConfigFileRepository_WithEmptyEnvVariable(t *testing.T) {
	// Test with empty environment variable
	originalEnv := os.Getenv("CONFIG_FILE_PATH")
	defer func() {
		if originalEnv != "" {
			os.Setenv("CONFIG_FILE_PATH", originalEnv)
		} else {
			os.Unsetenv("CONFIG_FILE_PATH")
		}
	}()

	// Clear the environment variable
	os.Unsetenv("CONFIG_FILE_PATH")

	// This should work with empty string (will create file in home directory root)
	repo := NewConfigFileRepository(os.Getenv("CONFIG_FILE_PATH"))

	// Verify it's a configFileRepository instance
	configRepo, ok := repo.(*configFileRepository)
	require.True(t, ok, "Should return configFileRepository instance")

	// Verify the file path is constructed correctly (empty string should result in just home directory)
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	expectedPath := filepath.Join(homeDir, "")
	assert.Equal(t, expectedPath, configRepo.filePath)

	// Note: This test demonstrates the behavior but doesn't actually save to avoid
	// creating files in the home directory during testing
}
