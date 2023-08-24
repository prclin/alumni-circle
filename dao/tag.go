package dao

import (
	"github.com/prclin/alumni-circle/model/po"
	"gorm.io/gorm"
)

type TagDao struct {
	Tx *gorm.DB
}

func NewTagDao(tx *gorm.DB) *TagDao {
	return &TagDao{Tx: tx}
}

func (td *TagDao) InsertBy(tag po.TTag) (uint32, error) {
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

func (td *TagDao) SelectById(id uint32) (po.TTag, error) {
	var tag po.TTag
	sql := "select id, name, state, extra, create_time, update_time from tag where id=?"
	err := td.Tx.Raw(sql, id).First(&tag).Error
	return tag, err
}

func (td *TagDao) UpdateTagBy(tag po.TTag) error {
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
