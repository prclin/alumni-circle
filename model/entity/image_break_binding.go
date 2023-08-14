package entity

import (
	"github.com/prclin/alumni-circle/global"
	"time"
)

type ImageBreakBinding struct {
	BreakId                 int       `json:"break_id" gorm:"column:break_id;primary_key"`
	ImageId                 int       `json:"image_id" gorm:"column:image_id" sql:"index"`
	Order                   int       `json:"order" gorm:"column:order"`
	CreateTime              time.Time `json:"create_time" gorm:"column:create_time;type:datetime(0);autoUpdateTime"`
	UpdateTime              time.Time `json:"update_time" gorm:"column:update_time;type:datetime(0);autoUpdateTime"`
	Deleted                 int       `json:"deleted" gorm:"column:deleted"`
	*ImageBreakBindingExtra `json:"extra" gorm:"column:extra"`
}
type ImageBreakBindingExtra struct{}

func (ImageBreakBinding) TableName() string {
	return "image_break_binding"
}

func CreateImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Create(binding).Error
}

func SaveImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Save(binding).Error
}

func GetImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Where(binding).Take(binding).Error
}

func CreateOrUpdateImageBreakBinding(binding *ImageBreakBinding) error {
	if GetImageBreakBinding(binding) == nil {
		return SaveImageBreakBinding(binding)
	} else {
		return CreateImageBreakBinding(binding)
	}
}
