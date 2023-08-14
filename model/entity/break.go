package entity

import (
	"github.com/prclin/alumni-circle/global"
	"time"
)

type Break struct {
	Id          int       `json:"id" gorm:"column:id;primary_key"`
	AccountId   int       `json:"account_id" gorm:"column:account_id"`
	Title       string    `json:"title" gorm:"column:title"`
	Content     string    `json:"content" gorm:"column:content"`
	Visibility  int       `json:"visibility" gorm:"column:visibility"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;type:datetime(0);autoUpdateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;type:datetime(0);autoUpdateTime"`
	Deleted     int       `json:"deleted" gorm:"column:deleted"`
	*BreakExtra `json:"extra" gorm:"column:extra"`
}
type BreakExtra struct{}

func (Break) TableName() string {
	return "break"
}

func CreateBreak(aBreak *Break) error {
	return global.Datasource.Create(aBreak).Error
}

func GetBreak(aBreak *Break) error {
	return global.Datasource.Where(aBreak).Take(aBreak).Error
}
