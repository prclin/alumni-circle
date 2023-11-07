package model

import "time"

// TComment 评论表
type TComment struct {
	Id         uint64    `json:"id"`
	ParentId   uint64    `json:"parent_id"`
	AccountId  uint64    `json:"account_id"`
	BreakId    uint64    `json:"break_id"`
	Content    string    `json:"content"`
	ReplyCount uint32    `json:"reply_count"`
	LikeCount  uint32    `json:"like_count"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// Comment 评论实体
type Comment struct {
	TComment
	Liked       bool         `json:"liked"`
	AccountInfo *AccountInfo `json:"account_info"`
}
