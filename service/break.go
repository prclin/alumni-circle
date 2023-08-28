package service

import (
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

func DeleteBreak(tBreak model.TBreak) error {
	breakDao := dao.NewBreakDao(global.Datasource)
	return breakDao.DeleteByIdAndAccountId(tBreak.Id, tBreak.AccountId)
}

func UpdateBreakVisibility(tBreak model.TBreak) error {
	bd := dao.NewBreakDao(global.Datasource)
	return bd.UpdateVisibilityBy(tBreak)
}

func PublishBreak(tBreak model.TBreak, shotIds, topicIds []uint64) (model.Break, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	bd := dao.NewBreakDao(tx)
	//创建课间
	breakId, err := bd.InsertBy(tBreak)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//绑定图片
	shotBindings := make([]model.TShotBinding, 0, len(shotIds))
	for index, shotId := range shotIds {
		shotBindings = append(shotBindings, model.TShotBinding{BreakId: breakId, ImageId: shotId, Order: uint8(index)})
	}
	sd := dao.NewShotDao(tx)
	err = sd.BatchInsertBy(shotBindings)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//绑定话题
	topicBindings := make([]model.TTopicBinding, 0, len(topicIds))
	for _, topicId := range topicIds {
		topicBindings = append(topicBindings, model.TTopicBinding{BreakId: breakId, TopicId: topicId})
	}
	td := dao.NewTopicDao(tx)
	err = td.BatchInsertBindingBy(topicBindings)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//获取课间
	var _break model.Break
	tb, err := bd.SelectById(breakId) //基本信息
	_break.TBreak = tb
	shots, err := sd.SelectShotsByBreakId(breakId) //镜头
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	_break.Shots = shots
	topics, err := td.SelectTopicsByBreakId(breakId) //话题
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	_break.Topics = topics

	return _break, nil
}
