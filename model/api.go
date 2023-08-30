package model

import "time"

// TAPI api表
type TAPI struct {
	Id          uint32
	Name        string
	Method      string
	Path        string
	Description string
	State       uint8
	Extra       *string
	CreateTime  time.Time
	UpdateTime  time.Time
}
