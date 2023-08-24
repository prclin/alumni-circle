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
	sql := "insert into tag(name) value (?)"
	//插入
	if err := td.Tx.Exec(sql, tag.Name).Error; err != nil {
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
	sql := "select id, name, extra, create_time, update_time from tag where id=?"
	err := td.Tx.Raw(sql, id).First(&tag).Error
	return tag, err
}

func (td *TagDao) UpdateTagBy(tag po.TTag) error {
	sql := "update tag set name=?,extra=? where id=?"
	return td.Tx.Exec(sql, tag.Name, tag.Extra, tag.Id).Error
}
