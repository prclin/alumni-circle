package entity

import "time"

type Account struct {
	Id         uint64
	Phone      string
	Email      string
	Password   string
	State      uint8
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}

type AccountInfo struct {
	Id         uint64
	CampusId   uint32
	AvatarURL  string
	Nickname   string
	Sex        uint8
	Birthday   time.Time
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}
