package postgres

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain/clusterDomain"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type configRepo struct {
	db *gorm.DB
}

type Config struct {
	gorm.Model

	ID             uuid.UUID
	LinkedObjectID uuid.UUID
	Name           string
	Value          string
}

func (c configRepo) GetConfigByObjectID(ID uuid.UUID) ([]clusterDomain.Environment, *e.Error) {
	var configs []Config
	if err := c.db.Where("linked_object_id = ?", ID).Find(&configs).Error; err != nil {
		return nil, e.Wrap(err)
	}

	return configsToDomain(configs), nil
}

func (c configRepo) SetConfig(ID uuid.UUID, config []clusterDomain.Environment) *e.Error {
	err := c.db.Where("linked_object_id = ?", ID).Delete(&Config{}).Error

	if err != nil {
		spew.Dump(err)
		return e.Wrap(err)
	}

	//TODO: replace with bulk insert and update if key already exists instead of delete and insert
	for _, env := range config {
		objectToSave := Config{
			ID:             uuid.New(),
			LinkedObjectID: ID,
			Name:           env.Name,
			Value:          env.Value,
		}

		if err := c.db.Save(&objectToSave).Error; err != nil {
			return e.Wrap(err)
		}
	}
	return nil
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
