package model

import "time"

// TBreak 课间表
type TBreak struct {
	Id         uint64
	AccountId  uint64
	Content    string
	Visibility uint8
	State      uint8
	Extra      *string
	CreateTime time.Time
	UpdateTime time.Time
}

type Break struct {
	TBreak
	Shots  []Shot
	Topics []TTopic
}
