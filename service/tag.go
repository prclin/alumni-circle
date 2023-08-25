package service

import (
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

func CreateTag(tag model.TTag) (model.TTag, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	td := dao.NewTagDao(tx)
	//插入
	id, err := td.InsertBy(tag)
	if err != nil {
		tx.Rollback()
		return tag, err
	}
	//获取标签
	tag, err = td.SelectById(id)
	if err != nil {
		tx.Rollback()
		return tag, err
	}
	return tag, nil
}

func UpdateTag(tag model.TTag) (model.TTag, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	td := dao.NewTagDao(tx)
	//更新
	err := td.UpdateTagBy(tag)
	if err != nil {
		tx.Rollback()
		return tag, err
	}
	//获取更新后tag
	tag, err = td.SelectById(tag.Id)
	if err != nil {
		tx.Rollback()
		return tag, err
	}
	return tag, nil
}

func DeleteTag(id uint32) error {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	td := dao.NewTagDao(tx)
	//删除关联的
	err := td.DeleteBindingByTagId(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	//删除tag
	err = td.DeleteById(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetTagList(pagination model.Pagination, state *uint8) ([]model.TTag, error) {
	td := dao.NewTagDao(global.Datasource)
	return td.SelectPageByState(state, (pagination.Page-1)*pagination.Size, pagination.Size)
}
