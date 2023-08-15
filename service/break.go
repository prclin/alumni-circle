package service

import (
	"github.com/prclin/alumni-circle/model/entity"
)

func CreateBreak(aBreak *entity.Break) (err error) {
	return entity.CreateBreak(aBreak)
}

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
