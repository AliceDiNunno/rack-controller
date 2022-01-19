package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type configRepo struct {
	db *gorm.DB
}

func (c configRepo) GetConfigByObjectID(ID uuid.UUID) ([]clusterDomain.Environment, *e.Error) {
	var configs []clusterDomain.Environment
	if err := c.db.Where("object_id = ?", ID).Find(&configs).Error; err != nil {
		return nil, e.Wrap(err)
	}
	return configs, nil
}

func (c configRepo) SetConfig(ID uuid.UUID, config []clusterDomain.Environment) *e.Error {
	for _, env := range config {
		if err := c.db.Save(&env).Error; err != nil {
			return e.Wrap(err)
		}
	}
	return nil
}

type Config struct {
	gorm.Model
	ClusterModel

	ID             uuid.UUID
	LinkedObjectID uuid.UUID
	Name           string
	Value          string
}

func configsToDomain(config []Config) []clusterDomain.Environment {
	configSlice := []clusterDomain.Environment{}

	for _, p := range config {
		configSlice = append(configSlice, configToDomain(p))
	}

	return configSlice
}

func configToDomain(config Config) clusterDomain.Environment {
	return clusterDomain.Environment{
		Name:  config.Name,
		Value: config.Value,
	}
}

func NewConfigRepo(db *gorm.DB) configRepo {
	return configRepo{
		db: db,
	}
}
