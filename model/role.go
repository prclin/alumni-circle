package model

import "time"

// TRoleBinding 角色绑定表
type TRoleBinding struct {
	AccountId  uint64
	RoleId     uint32
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}
