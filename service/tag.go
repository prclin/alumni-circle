package service

import (
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/po"
)

func CreateTag(tag po.TTag) (po.TTag, error) {
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

func UpdateTag(tag po.TTag) (po.TTag, error) {
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
