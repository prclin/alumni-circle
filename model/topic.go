package model

import "time"

// TTopic 话题表
type TTopic struct {
	Id         uint64    `json:"id"`
	Name       string    `json:"name"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// TTopicBinding 话题绑定表
type TTopicBinding struct {
	BreakId uint64
	TopicId uint64
}
