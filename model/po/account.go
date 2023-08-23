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
	Id            uint64     `json:"id"`
	CampusId      uint32     `json:"campus_id"`
	AvatarURL     string     `json:"avatar_url"`
	Nickname      string     `json:"nickname"`
	Sex           uint8      `json:"sex"`
	Birthday      *time.Time `json:"birthday"`
	FollowCount   uint32     `json:"follow_count"`
	FollowerCount uint32     `json:"follower_count"`
	Extra         string     `json:"extra"`
	CreateTime    time.Time  `json:"create_time"`
	UpdateTime    time.Time  `json:"update_time"`
}
