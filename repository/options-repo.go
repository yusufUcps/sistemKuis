package repository

import (
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OptionsInterface interface {
	InsertOption(newOptions model.Options) (*model.Options, int)
}

type OptionsModel struct {
	db *gorm.DB
}

func (qm *OptionsModel) Init(db *gorm.DB) {
	qm.db = db
}

func NewOptionsModel(db *gorm.DB) OptionsInterface {
	return &OptionsModel{
		db: db,
	}
}

func (om *OptionsModel) InsertOption(newOptions model.Options) (*model.Options, int) {

	if err := om.db.Create(&newOptions).Error; err != nil {
		logrus.Error("Repository: Insert data Option error, ", err.Error())
		return nil, 1
	}

	return &newOptions  , 0
}