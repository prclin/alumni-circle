package service

import (
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
)

func AddImageToBreak(breakId int, bindingList []entity.ImageBreakBinding) (err error) {
	for _, binding := range bindingList {
		binding.BreakId = breakId
		if err = entity.CreateOrUpdateImageBreakBinding(&binding); err != nil {
			return err
		}
	}
	sizeLimit := global.Configuration.Limit.Size.PictureInBreak
	for order := len(bindingList); order < sizeLimit; order++ {
		binding := &entity.ImageBreakBinding{
			BreakId: breakId,
			Order:   order,
		}
		if entity.GetImageBreakBinding(binding) != nil {
			if err := entity.DeleteImageBreakBinding(binding); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

func ImageBreakBindingExist(breakId int, order int) bool {
	binding := &entity.ImageBreakBinding{
		BreakId: breakId,
		Order:   order,
	}
	if err := entity.GetImageBreakBinding(binding); err != nil {
		return false
	}
	return true
}
