package config

import (
	"errors"
	"slices"

	"github.com/kyoshidaxx/tsunagi/internal/utils"
)

type ConfigParam struct {
	Name         string
	Port         int
	ProjectName  string
	Region       string
	InstanceName string
}

type Config struct {
	r Repository
}

const (
	ephemelalPortFrom = 49152
	ephemelalPortTo   = 65535
)

func NewConfig(r Repository) *Config {
	return &Config{r: r}
}

func (c *Config) Add(param ConfigParam) error {
	if len(param.Name) == 0 {
		return errors.New("name is required")
	}
	if param.Port < ephemelalPortFrom || param.Port > ephemelalPortTo {
		return errors.New("port is out of range")
	}
	if len(param.ProjectName) == 0 {
		return errors.New("project name is required")
	}
	if len(param.Region) == 0 {
		return errors.New("region is required")
	}
	regionList := utils.GetRegionList()
	if !slices.Contains(regionList, param.Region) {
		return errors.New("region is not valid")
	}
	if len(param.InstanceName) == 0 {
		return errors.New("instance name is required")
	}
	// todo Nameの重複チェック

	return c.r.Save(param)
}
