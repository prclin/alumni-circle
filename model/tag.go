package model

import "time"

// TTag 标签表
type TTag struct {
	Id         uint32    `json:"id"`
	Name       string    `json:"name"`
	State      uint8     `json:"state"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// TAccountTagBinding 用户标签绑定表
type TAccountTagBinding struct {
	AccountId uint64
	TagId     uint32
}

// TBreakTagBinding 课间标签绑定表
type TBreakTagBinding struct {
	BreakId uint64
	TagId   uint32
}
