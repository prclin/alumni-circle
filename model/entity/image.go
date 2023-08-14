package entity

import (
	"github.com/prclin/alumni-circle/global"
	"time"
)

type Image struct {
	Id          int       `json:"id" gorm:"column:id;primary_key"`
	Url         string    `json:"url" gorm:"column:url"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;type:datetime(0);autoUpdateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;type:datetime(0);autoUpdateTime"`
	Deleted     int       `json:"deleted" gorm:"column:deleted"`
	*ImageExtra `json:"extra" gorm:"column:extra"`
}
type ImageExtra struct{}

func (Image) TableName() string {
	return "image"
}

func CreateImage(image *Image) error {
	return global.Datasource.Create(&image).Error
}
