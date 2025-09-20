package datastore

import (
	config "github.com/kyoshidaxx/tsunagi/internal/domain/config"
)

type configFileRepository struct {
}

func NewConfigFileRepository() config.Repository {
	return &configFileRepository{}
}

func (r *configFileRepository) Save(config config.ConfigParam) error {
	// todo
	return nil
}
