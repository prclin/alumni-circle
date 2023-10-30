package model

import (
	"strconv"
	"time"
)

// TBreak 课间表
type TBreak struct {
	Id         uint64
	AccountId  uint64
	Content    string
	Visibility uint8
	State      uint8
	LikeCount  uint32
	Extra      *string
	CreateTime time.Time
	UpdateTime time.Time
}

type Break struct {
	*TBreak
	Shots  []Shot
	Topics []TTopic
	Tags   []TTag
}

// TBreakLike 课件点赞表
type TBreakLike struct {
	AccountId uint64
	BreakId   uint64
}

func (tbl *TBreakLike) String() string {
	return strconv.FormatUint(tbl.AccountId, 10) + ":" + strconv.FormatUint(tbl.BreakId, 10)
}
