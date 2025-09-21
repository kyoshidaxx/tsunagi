package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepository is a mock implementation of the Repository interface for testing
type mockRepository struct {
	saveCalled bool
	saveParam  ConfigParam
	saveError  error
}

func (m *mockRepository) Save(config ConfigParam) error {
	m.saveCalled = true
	m.saveParam = config
	return m.saveError
}

func (m *mockRepository) reset() {
	m.saveCalled = false
	m.saveParam = ConfigParam{}
	m.saveError = nil
}

func TestNewConfig(t *testing.T) {
	mockRepo := &mockRepository{}
	config := NewConfig(mockRepo)

	assert.NotNil(t, config)
	assert.Equal(t, mockRepo, config.r)
}

func TestConfig_Add_Success(t *testing.T) {
	mockRepo := &mockRepository{}
	config := NewConfig(mockRepo)

	param := ConfigParam{
		Name:         "test-config",
		Port:         50000,
		ProjectName:  "test-project",
		Region:       "asia-northeast1",
		InstanceName: "test-instance",
	}

	err := config.Add(param)

	assert.NoError(t, err)
	assert.True(t, mockRepo.saveCalled)
	assert.Equal(t, param, mockRepo.saveParam)
}

func TestConfig_Add_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		param   ConfigParam
		wantErr string
	}{
		{
			name: "empty name",
			param: ConfigParam{
				Name:         "",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "name is required",
		},
		{
			name: "port too low",
			param: ConfigParam{
				Name:         "test-config",
				Port:         1000,
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "port is out of range",
		},
		{
			name: "port too high",
			param: ConfigParam{
				Name:         "test-config",
				Port:         70000,
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "port is out of range",
		},
		{
			name: "port at minimum boundary",
			param: ConfigParam{
				Name:         "test-config",
				Port:         49151, // 49152未満
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "port is out of range",
		},
		{
			name: "port at maximum boundary",
			param: ConfigParam{
				Name:         "test-config",
				Port:         65536, // 65535超過
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "port is out of range",
		},
		{
			name: "empty project name",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
			wantErr: "project name is required",
		},
		{
			name: "empty region",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "",
				InstanceName: "test-instance",
			},
			wantErr: "region is required",
		},
		{
			name: "invalid region",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "invalid-region",
				InstanceName: "test-instance",
			},
			wantErr: "region is not valid",
		},
		{
			name: "empty instance name",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "",
			},
			wantErr: "instance name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			config := NewConfig(mockRepo)

			err := config.Add(tt.param)

			assert.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
			assert.False(t, mockRepo.saveCalled, "Save should not be called when validation fails")
		})
	}
}

func TestConfig_Add_ValidBoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		param ConfigParam
	}{
		{
			name: "minimum valid port",
			param: ConfigParam{
				Name:         "test-config",
				Port:         49152, // min
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
		},
		{
			name: "maximum valid port",
			param: ConfigParam{
				Name:         "test-config",
				Port:         65535, // max
				ProjectName:  "test-project",
				Region:       "asia-northeast1",
				InstanceName: "test-instance",
			},
		},
		{
			name: "valid region asia-east1",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "asia-east1",
				InstanceName: "test-instance",
			},
		},
		{
			name: "valid region us-west4",
			param: ConfigParam{
				Name:         "test-config",
				Port:         50000,
				ProjectName:  "test-project",
				Region:       "us-west4",
				InstanceName: "test-instance",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			config := NewConfig(mockRepo)

			err := config.Add(tt.param)

			assert.NoError(t, err)
			assert.True(t, mockRepo.saveCalled)
			assert.Equal(t, tt.param, mockRepo.saveParam)
		})
	}
}

func TestConfig_Add_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{}
	config := NewConfig(mockRepo)

	// simulate repository save failed
	expectedError := errors.New("repository save failed")
	mockRepo.saveError = expectedError

	param := ConfigParam{
		Name:         "test-config",
		Port:         50000,
		ProjectName:  "test-project",
		Region:       "asia-northeast1",
		InstanceName: "test-instance",
	}

	err := config.Add(param)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.True(t, mockRepo.saveCalled)
	assert.Equal(t, param, mockRepo.saveParam)
}

func TestConfig_Add_MultipleValidationErrors(t *testing.T) {
	// multiple validation errors, first error is returned
	mockRepo := &mockRepository{}
	config := NewConfig(mockRepo)

	param := ConfigParam{
		Name:         "",   // empty name
		Port:         1000, // out of range port
		ProjectName:  "",   // empty project name
		Region:       "",   // empty region
		InstanceName: "",   // empty instance name
	}

	err := config.Add(param)

	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
	assert.False(t, mockRepo.saveCalled)
}

// benchmark test
func BenchmarkConfig_Add(b *testing.B) {
	mockRepo := &mockRepository{}
	config := NewConfig(mockRepo)

	param := ConfigParam{
		Name:         "benchmark-config",
		Port:         50000,
		ProjectName:  "benchmark-project",
		Region:       "asia-northeast1",
		InstanceName: "benchmark-instance",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockRepo.reset()
		err := config.Add(param)
		require.NoError(b, err)
	}
}
