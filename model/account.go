package model

import (
	"time"
)

// Account 数据实体
type Account struct {
	Info       TAccountInfo `json:"info"`
	IsFollowed bool         `json:"is_followed"`
}

// TAccount 账户表
type TAccount struct {
	Id         uint64
	Phone      string
	Email      string
	Password   string `json:"-"`
	State      uint8
	Extra      *string
	CreateTime time.Time
	UpdateTime time.Time
}

// TAccountInfo 账户信息表
type TAccountInfo struct {
	Id            uint64    `json:"id"`
	CampusId      uint32    `json:"campus_id"`
	AvatarURL     string    `json:"avatar_url"`
	Nickname      string    `json:"nickname"`
	Sex           uint8     `json:"sex"`
	Birthday      time.Time `json:"birthday"`
	FollowCount   uint32    `json:"follow_count"`
	FollowerCount uint32    `json:"follower_count"`
	Extra         *string   `json:"extra"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}
