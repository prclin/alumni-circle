package service

import (
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

func CreateTopic(topic model.TTopic) (model.TTopic, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	td := dao.NewTopicDao(tx)
	//插入话题
	id, err := td.InsertBy(topic)
	if err != nil {
		tx.Rollback()
		return model.TTopic{}, err
	}
	//获取话题
	tTopic, err := td.SelectById(id)
	if err != nil {
		tx.Rollback()
		return model.TTopic{}, err
	}
	return tTopic, nil
}
