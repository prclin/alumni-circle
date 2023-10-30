package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strconv"
	"strings"
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

func (bd *BreakDao) SelectById(id uint64) (*model.TBreak, error) {
	var tBreak *model.TBreak
	sql := "select id, account_id, content, visibility, state, like_count, extra, create_time, update_time from break where id=?"
	err := bd.Tx.Raw(sql, id).First(&tBreak).Error
	return tBreak, err
}

func (bd *BreakDao) SelectByIds(ids []uint64) ([]model.TBreak, error) {
	var breaks []model.TBreak
	sql := "select id, account_id, content, visibility, state, like_count, extra, create_time, update_time from break where id in ?"
	err := bd.Tx.Raw(sql, ids).Scan(&breaks).Error
	return breaks, err
}

func (bd *BreakDao) UpdateVisibilityBy(tBreak model.TBreak) error {
	sql := "update break set visibility=? where id=? and account_id=?"
	return bd.Tx.Exec(sql, tBreak.Visibility, tBreak.Id, tBreak.AccountId).Error
}

func (bd *BreakDao) DeleteByIdAndAccountId(id, accountId uint64) error {
	sql := "delete from break where id=? and account_id=?"
	return bd.Tx.Exec(sql, id, accountId).Error
}

func (bd *BreakDao) SelectApprovedIdsRandomlyBefore(latestTime int64, accountId uint64, limit int) ([]uint64, error) {
	var ids []uint64
	sql := "select id from break where account_id != ? and state=2 and create_time >= ? order by RAND() limit ?"
	err := bd.Tx.Raw(sql, accountId, latestTime, limit).Scan(&ids).Error
	return ids, err
}

func (bd *BreakDao) BatchInsertLikeBy(likes []model.TBreakLike) error {
	if likes == nil || len(likes) == 0 {
		return nil
	}
	var sql strings.Builder
	sql.WriteString("insert into break_like values ")
	params := make([]interface{}, 0, len(likes)*2)
	for _, like := range likes {
		sql.WriteString("(?,?),")
		params = append(params, like.AccountId, like.BreakId)
	}
	return bd.Tx.Exec(sql.String()[:sql.Len()-1], params).Error
}

func (bd *BreakDao) BatchDeleteLikeBy(unlikes []model.TBreakLike) error {
	if unlikes == nil || len(unlikes) == 0 {
		return nil
	}

	sql := "delete from break_like where (account_id,break_id) in (?)"
	var param strings.Builder
	for _, unlike := range unlikes {
		param.WriteString("(")
		param.WriteString(strconv.FormatUint(unlike.AccountId, 10))
		param.WriteString(",")
		param.WriteString(strconv.FormatUint(unlike.BreakId, 10))
		param.WriteString("),")
	}
	return bd.Tx.Exec(sql, param.String()[:param.Len()-1]).Error
}

func (bd *BreakDao) BatchIncreaseLikeCount(increases map[uint64]uint32) error {
	pattern := "update break set like_count=like_count + ? where id = ?;"
	var sql strings.Builder
	params := make([]any, 0, len(increases)*2)
	for key, value := range increases {
		sql.WriteString(pattern)
		params = append(params, key, value)
	}
	return bd.Tx.Exec(sql.String(), params).Error
}
