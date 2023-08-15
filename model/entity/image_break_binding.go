package entity

import (
	"github.com/prclin/alumni-circle/global"
	"time"
)

type ImageBreakBinding struct {
	BreakId                 int       `json:"break_id" gorm:"column:break_id;primary_key"`
	ImageId                 int       `json:"image_id" gorm:"column:image_id" sql:"index"`
	Order                   int       `json:"order" gorm:"column:order;primary_key"`
	CreateTime              time.Time `json:"create_time" gorm:"column:create_time;type:datetime(0);autoUpdateTime"`
	UpdateTime              time.Time `json:"update_time" gorm:"column:update_time;type:datetime(0);autoUpdateTime"`
	Deleted                 int       `json:"deleted" gorm:"column:deleted"`
	*ImageBreakBindingExtra `json:"extra" gorm:"column:extra"`
}
type ImageBreakBindingExtra struct{}

func (ImageBreakBinding) TableName() string {
	return "image_break_binding"
}

func GetImageBreakBinding(binding *ImageBreakBinding) error {
	tx := global.Datasource.Where("deleted=0")
	tx = tx.Where("break_id=?", binding.BreakId)
	tx = tx.Where("order=?", binding.Order)
	return tx.Take(binding).Error
}

func CreateImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Create(binding).Error
}

func UpdateImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Updates(binding).Error
}

func DeleteImageBreakBinding(binding *ImageBreakBinding) error {
	return global.Datasource.Delete(binding).Error
}

func CreateOrUpdateImageBreakBinding(binding *ImageBreakBinding) error {
	if err := GetImageBreakBinding(
		&ImageBreakBinding{
			BreakId: binding.BreakId,
			Order:   binding.Order,
		},
	); err != nil {
		return CreateImageBreakBinding(binding)
	} else {
		return UpdateImageBreakBinding(binding)
	}
}
