package model

import "time"

// TRole 角色表
type TRole struct {
	Id          uint32    `json:"id"`
	Name        string    `json:"name"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	State       uint8     `json:"state"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// TRoleBinding 角色绑定表
type TRoleBinding struct {
	AccountId  uint64
	RoleId     uint32
	Extra      *string
	CreateTime time.Time
	UpdateTime time.Time
}
