package datastore

import (
	"encoding/json"
	"os"
	"path/filepath"

	c "github.com/kyoshidaxx/tsunagi/internal/domain/config"
)

type configFileRepository struct {
	filePath string
}

func NewConfigFileRepository(filePath string) c.Repository {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filePath = filepath.Join(homeDir, filePath)
	return &configFileRepository{filePath: filePath}
}

func (r *configFileRepository) Save(config c.ConfigParam) error {
	fileExists := r.checkConfigFileExists()
	var configParams []c.ConfigParam

	if fileExists {
		var err error
		configParams, err = r.loadConfigFile()
		if err != nil {
			return err
		}
	} else {
		err := r.createConfigFile()
		if err != nil {
			return err
		}
	}

	configParams = append(configParams, config)
	data, err := json.MarshalIndent(configParams, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(r.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *configFileRepository) checkConfigFileExists() bool {
	_, err := os.Stat(r.filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (r *configFileRepository) loadConfigFile() ([]c.ConfigParam, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []c.ConfigParam{}, nil
	}

	var configParams []c.ConfigParam
	err = json.Unmarshal(data, &configParams)
	if err != nil {
		return nil, err
	}

	return configParams, nil
}

func (r *configFileRepository) createConfigFile() error {

	dir := filepath.Dir(r.filePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	_, err = os.Create(r.filePath)
	if err != nil {
		return err
	}
	return nil
}
