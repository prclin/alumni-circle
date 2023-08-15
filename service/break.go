package service

import (
	"github.com/prclin/alumni-circle/model/entity"
)

// CreateBreak 创建课间
func CreateBreak(aBreak *entity.Break) (err error) {
	return entity.CreateBreak(aBreak)
}

// BreakExist 课间是否存在
func BreakExist(breakId int, accountId int) bool {
	aBreak := &entity.Break{
		Id:        breakId,
		AccountId: accountId,
	}
	if err := entity.GetBreak(aBreak); err != nil {
		return false
	}
	return true
}
