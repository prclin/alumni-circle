package model

import "time"

// TAPI api表
type TAPI struct {
	Id          uint32    `json:"id"`
	Name        string    `json:"name"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	State       uint8     `json:"state"`
	Extra       *string   `json:"extra"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// TAPIBinding api_binding表
type TAPIBinding struct {
	RoleId uint32 `json:"role_id"`
	APIId  uint32 `json:"api_id"`
}
