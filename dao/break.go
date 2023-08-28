package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
)

type BreakDao struct {
	Tx *gorm.DB
}

func NewBreakDao(tx *gorm.DB) *BreakDao {
	return &BreakDao{Tx: tx}
}

func (bd *BreakDao) InsertBy(tBreak model.TBreak) (uint64, error) {
	var id uint64
	sql := "insert into break(account_id, content, visibility, state, extra) value (?,?,?,?,?)"
	//插入数据
	if err := bd.Tx.Exec(sql, tBreak.AccountId, tBreak.Content, tBreak.Visibility, tBreak.State, tBreak.Extra).Error; err != nil {
		return 0, err
	}
	//获取id
	if err := bd.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (bd *BreakDao) SelectById(id uint64) (model.TBreak, error) {
	var tBreak model.TBreak
	sql := "select id, account_id, content, visibility, state, extra, create_time, update_time from break where id=?"
	err := bd.Tx.Raw(sql, id).First(&tBreak).Error
	return tBreak, err
}

func (bd *BreakDao) UpdateVisibilityBy(tBreak model.TBreak) error {
	sql := "update break set visibility=? where id=? and account_id=?"
	return bd.Tx.Exec(sql, tBreak.Visibility, tBreak.Id, tBreak.AccountId).Error
}
