package service

import (
	"github.com/prclin/alumni-circle/model/entity"
)

func CreateBreak(aBreak *entity.Break) (err error) {
	return entity.CreateBreak(aBreak)
}

func BreakExist(breakId int, accountId int) error {
	aBreak := &entity.Break{
		Id:        breakId,
		AccountId: accountId,
	}
	return entity.GetBreak(aBreak)
}

func AddImageToBreak(binding *entity.ImageBreakBinding) error {
	return entity.CreateOrUpdateImageBreakBinding(binding)
}
