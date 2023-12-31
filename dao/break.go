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

func (dao *BreakDao) InsertBy(tBreak model.TBreak) (uint64, error) {
	var id uint64
	sql := "insert into break(account_id, content, visibility, state, extra) value (?,?,?,?,?)"
	//插入数据
	if err := dao.Tx.Exec(sql, tBreak.AccountId, tBreak.Content, tBreak.Visibility, tBreak.State, tBreak.Extra).Error; err != nil {
		return 0, err
	}
	//获取id
	if err := dao.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (dao *BreakDao) SelectById(id uint64) (*model.TBreak, error) {
	var tBreak *model.TBreak
	sql := "select id, account_id, content, visibility,like_count, comment_count, state,  extra, create_time, update_time from break where id=?"
	err := dao.Tx.Raw(sql, id).First(&tBreak).Error
	return tBreak, err
}

func (dao *BreakDao) SelectByIds(ids []uint64) ([]model.TBreak, error) {
	var breaks []model.TBreak
	sql := "select id, account_id, content, visibility,like_count, comment_count, state, extra, create_time, update_time from break where id in ?"
	err := dao.Tx.Raw(sql, ids).Scan(&breaks).Error
	return breaks, err
}

func (dao *BreakDao) UpdateVisibilityBy(tBreak model.TBreak) error {
	sql := "update break set visibility=? where id=? and account_id=?"
	return dao.Tx.Exec(sql, tBreak.Visibility, tBreak.Id, tBreak.AccountId).Error
}

func (dao *BreakDao) DeleteByIdAndAccountId(id, accountId uint64) error {
	sql := "delete from break where id=? and account_id=?"
	return dao.Tx.Exec(sql, id, accountId).Error
}

func (dao *BreakDao) SelectApprovedIdsRandomlyBefore(latestTime int64, accountId uint64, limit int) ([]uint64, error) {
	var ids []uint64
	sql := "select id from break where account_id != ? and state=2 and create_time >= ? order by RAND() limit ?"
	err := dao.Tx.Raw(sql, accountId, latestTime, limit).Scan(&ids).Error
	return ids, err
}

func (dao *BreakDao) BatchInsertLikeBy(likes []model.TBreakLike) error {
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
	return dao.Tx.Exec(sql.String()[:sql.Len()-1], params).Error
}

func (dao *BreakDao) BatchDeleteLikeBy(unlikes []model.TBreakLike) error {
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
	return dao.Tx.Exec(sql, param.String()[:param.Len()-1]).Error
}

func (dao *BreakDao) BatchIncreaseLikeCount(increases map[uint64]uint32) error {
	pattern := "update break set like_count=like_count + ? where id = ?;"
	var sql strings.Builder
	params := make([]any, 0, len(increases)*2)
	for key, value := range increases {
		sql.WriteString(pattern)
		params = append(params, key, value)
	}
	return dao.Tx.Exec(sql.String(), params).Error
}

func (dao *BreakDao) SelectByAccountIdAndVisibility(accountId uint64, visibility uint8, pagination model.Pagination) ([]model.TBreak, error) {
	var breaks []model.TBreak
	sql := "select * from break where account_id = ? and visibility >= ? limit ?,?"
	err := dao.Tx.Raw(sql, accountId, visibility, (pagination.Page-1)*pagination.Size, pagination.Size).Scan(&breaks).Error
	return breaks, err
}

func (dao *BreakDao) IsLiked(accountId, breakId uint64) bool {
	var liked bool
	sql := "select count(*) from break_like where  account_id = ? and break_id = ?"
	dao.Tx.Raw(sql, accountId, breakId).First(&liked)
	return liked
}
