package entity

import (
	"github.com/prclin/alumni-circle/global"
	"time"
)

type Image struct {
	Id          int       `json:"id" gorm:"primary_key"`
	Url         string    `json:"url"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	Deleted     int       `json:"deleted"`
	*ImageExtra `json:"extra"`
}
type ImageExtra struct{}

func (Image) TableName() string {
	return "image"
}
func CreateImage(image *Image) error {
	image.CreateTime = time.Now()
	image.UpdateTime = time.Now()
	return global.Datasource.Create(&image).Error
}
