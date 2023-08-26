package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
)

type TagDao struct {
	Tx *gorm.DB
}

func NewTagDao(tx *gorm.DB) *TagDao {
	return &TagDao{Tx: tx}
}

func (td *TagDao) InsertBy(tag model.TTag) (uint32, error) {
	var id uint32
	sql := "insert into tag(name, state, extra) value (?,?,?)"
	//插入
	if err := td.Tx.Exec(sql, tag.Name, tag.State, tag.Extra).Error; err != nil {
		return 0, err
	}
	//查询主键
	if err := td.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil

}

func (td *TagDao) SelectById(id uint32) (model.TTag, error) {
	var tag model.TTag
	sql := "select id, name, state, extra, create_time, update_time from tag where id=?"
	err := td.Tx.Raw(sql, id).First(&tag).Error
	return tag, err
}

func (td *TagDao) UpdateTagBy(tag model.TTag) error {
	sql := "update tag set name=?,state=?,extra=? where id=?"
	return td.Tx.Exec(sql, tag.Name, tag.State, tag.Extra, tag.Id).Error
}

func (td *TagDao) DeleteBindingByTagId(tagId uint32) error {
	sql := "delete from tag_binding where tag_id=?"
	return td.Tx.Exec(sql, tagId).Error
}

func (td *TagDao) DeleteById(id uint32) error {
	sql := "delete from tag where id=?"
	return td.Tx.Exec(sql, id).Error
}

func (td *TagDao) SelectPageByState(state *uint8, offset int, size int) ([]model.TTag, error) {
	var tags []model.TTag
	sql := "select id, name, state, extra, create_time, update_time from tag "
	params := make([]interface{}, 0, 1)
	if state != nil {
		sql += "where state=? "
		params = append(params, state)
	}
	sql += "limit ?,?"
	params = append(params, offset, size)
	err := td.Tx.Raw(sql, params...).Scan(&tags).Error
	return tags, err
}

func (td *TagDao) DeleteBindingByAccountId(accountId uint64) error {
	sql := "delete from tag_binding where account_id=?"
	return td.Tx.Exec(sql, accountId).Error
}

func (td *TagDao) BatchInsertBindingBy(bindings []model.TTagBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into tag_binding(account_id, tag_id) VALUES " //goland报错，忽略
	params := make([]interface{}, 0, 0)
	for _, binding := range bindings {
		sql += "(?,?),"
		params = append(params, binding.AccountId, binding.TagId)
	}
	sql = strings.TrimSuffix(sql, ",")
	return td.Tx.Exec(sql, params...).Error
}

func (td *TagDao) SelectEnabledByIds(ids []uint32) ([]model.TTag, error) {
	var tags []model.TTag
	sql := "select id, name, state, extra, create_time, update_time from tag where id in ? and state=1"
	err := td.Tx.Raw(sql, ids).Scan(&tags).Error
	return tags, err
}
