package po

import "time"

type TRoleBinding struct {
	AccountId  uint64
	RoleId     uint32
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}
