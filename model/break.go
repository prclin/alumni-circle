package model

import (
	"strconv"
	"time"
)

// TBreak 课间表
type TBreak struct {
	Id         uint64    `json:"id"`
	AccountId  uint64    `json:"account_id"`
	Content    string    `json:"content"`
	Visibility uint8     `json:"visibility"`
	State      uint8     `json:"state"`
	LikeCount  uint32    `json:"like_count"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type Break struct {
	TBreak
	Shots []Shot `json:"shots"`
	Tags  []TTag `json:"tags"`
}

// TBreakLike 课件点赞表
type TBreakLike struct {
	AccountId uint64
	BreakId   uint64
}

func (tbl *TBreakLike) String() string {
	return strconv.FormatUint(tbl.AccountId, 10) + ":" + strconv.FormatUint(tbl.BreakId, 10)
}
