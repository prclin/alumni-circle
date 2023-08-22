package po

import "time"

type TAccount struct {
	Id         uint64
	Phone      string
	Email      string
	Password   string `json:"-"`
	State      uint8
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}

type TAccountInfo struct {
	Id         uint64    `json:"id"`
	CampusId   uint32    `json:"campus_id"`
	AvatarURL  string    `json:"avatar_url"`
	Nickname   string    `json:"nickname"`
	Sex        uint8     `json:"sex"`
	Birthday   time.Time `json:"birthday"`
	Extra      string    `json:"-"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}
