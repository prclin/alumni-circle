package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
)

type TopicDao struct {
	Tx *gorm.DB
}

func NewTopicDao(tx *gorm.DB) *TopicDao {
	return &TopicDao{Tx: tx}
}

func (td *TopicDao) BatchInsertBindingBy(bindings []model.TTopicBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into topic_binding(break_id, topic_id) values" //此处为goland报错
	params := make([]interface{}, 0, len(bindings))
	for _, binding := range bindings {
		sql += "(?,?),"
		params = append(params, binding.BreakId, binding.TopicId)
	}
	sql = strings.TrimSuffix(sql, ",")
	return td.Tx.Exec(sql, params...).Error
}

func (td *TopicDao) SelectTopicsByBreakId(breakId uint64) ([]model.TTopic, error) {
	var topics []model.TTopic
	sql := "select id, name, extra, create_time, update_time from topic where id in (select topic_id from topic_binding where break_id=?)"
	err := td.Tx.Raw(sql, breakId).Scan(&topics).Error
	if topics == nil {
		topics = make([]model.TTopic, 0, 0)
	}
	return topics, err
}

func (td *TopicDao) InsertBy(topic model.TTopic) (uint64, error) {
	var id uint64
	sql := "insert into topic(name, extra) value (?,?)"
	//插入数据
	if err := td.Tx.Exec(sql, topic.Name, topic.Extra).Error; err != nil {
		return 0, err
	}
	//查询主键
	if err := td.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (td *TopicDao) SelectById(id uint64) (model.TTopic, error) {
	var topic model.TTopic
	sql := "select id, name, extra, create_time, update_time from topic where id=?"
	err := td.Tx.Raw(sql, id).First(&topic).Error
	return topic, err
}
