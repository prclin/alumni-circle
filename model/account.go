package model

import (
	"time"
)

// Account 数据实体
type Account struct {
	TAccount
	Info *AccountInfo `json:"info"`
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

type AccountInfo struct {
	TAccountInfo
	Tags           []TTag   `json:"tags"`
	Campus         *TCampus `json:"campus"`
	Followed       bool     `json:"is_followed"`
	MutualFollowed bool     `json:"mutual_followed"`
}

// TAccountInfo 账户信息表
type TAccountInfo struct {
	Id            uint64     `json:"id"`
	CampusId      uint32     `json:"campus_id"`
	Nickname      string     `json:"nickname"`
	AvatarURL     string     `json:"avatar_url"`
	BackgroundURL string     `json:"background_url"`
	Sex           uint8      `json:"sex"`
	Brief         string     `json:"brief"`
	Birthday      *time.Time `json:"birthday"`
	MBTIResultID  *string    `json:"mbti_result_id"`
	FollowCount   uint32     `json:"follow_count"`
	FollowerCount uint32     `json:"follower_count"`
	FriendCount   uint32     `json:"friend_count"`
	Extra         *string    `json:"extra"`
	CreateTime    time.Time  `json:"create_time"`
	UpdateTime    time.Time  `json:"update_time"`
}

// TFollow 关注表
type TFollow struct {
	FollowerId uint64
	FolloweeId uint64
	Extra      *string
	CreateTime time.Time
}
