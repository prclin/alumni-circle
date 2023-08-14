package service

import "github.com/prclin/alumni-circle/model/entity"

func CreateBreak(aBreak *entity.Break) (err error) {
	return entity.CreateBreak(aBreak)
}
