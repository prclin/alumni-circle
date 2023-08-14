package entity

import (
	"time"
)

type Topic struct {
	Id          int       `json:"id" gorm:"column:id;primary_key"`
	Name        string    `json:"name" gorm:"column:name"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;type:datetime(0);autoUpdateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;type:datetime(0);autoUpdateTime"`
	Deleted     int       `json:"deleted" gorm:"column:deleted"`
	*TopicExtra `json:"extra" gorm:"column:extra"`
}
type TopicExtra struct{}

func (Topic) TableName() string {
	return "image"
}
