package config

type Repository interface {
	Save(config ConfigParam) error
}
