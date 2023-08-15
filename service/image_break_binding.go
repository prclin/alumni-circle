package service

import (
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
)

// 添加图片至课间
func AddImageToBreak(breakId int, bindingList []entity.ImageBreakBinding) (err error) {
	// 遍历列表,将传入的图片添加至课间
	for _, binding := range bindingList {
		binding.BreakId = breakId
		if err = entity.CreateOrUpdateImageBreakBinding(&binding); err != nil {
			return err
		}
	}
	// 删除以往图片
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
