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
	sql.WriteString("insert ignore into break_like(account_id, break_id) values ")
	params := make([]interface{}, 0, len(likes)*2)
	for _, like := range likes {
		sql.WriteString("(?,?),")
		params = append(params, like.AccountId, like.BreakId)
	}
	return dao.Tx.Exec(sql.String()[:sql.Len()-1], params...).Error
}

func (dao *BreakDao) BatchDeleteLikeBy(unlikes []model.TBreakLike) error {
	length := len(unlikes)
	if unlikes == nil || length == 0 {
		return nil
	}

	var sql strings.Builder
	sql.WriteString("delete from break_like where (account_id,break_id) in ( ")
	for i := 0; i < length; i++ {
		sql.WriteString("(")
		sql.WriteString(strconv.FormatUint(unlikes[i].AccountId, 10))
		sql.WriteString(",")
		sql.WriteString(strconv.FormatUint(unlikes[i].BreakId, 10))
		sql.WriteString(")")

		if i == length-1 {
			sql.WriteString(")")
		} else {
			sql.WriteString(",")
		}
	}
	return dao.Tx.Exec(sql.String()).Error
}

func (dao *BreakDao) BatchIncreaseLikeCount(increases map[uint64]int) error {
	sql := "update break set like_count=like_count + ? where id = ?;"
	for key, value := range increases {
		err := dao.Tx.Exec(sql, value, key).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (dao *BreakDao) SelectByAccountIdAndVisibility(accountId uint64, visibility uint8, pagination model.Pagination) ([]model.TBreak, error) {
	var breaks []model.TBreak
	sql := "select * from break where account_id = ? and visibility >= ? order by create_time desc limit ?,?"
	err := dao.Tx.Raw(sql, accountId, visibility, (pagination.Page-1)*pagination.Size, pagination.Size).Scan(&breaks).Error
	return breaks, err
}

func (dao *BreakDao) IsLiked(accountId, breakId uint64) bool {
	var liked bool
	sql := "select count(*) from break_like where  account_id = ? and break_id = ?"
	dao.Tx.Raw(sql, accountId, breakId).First(&liked)
	return liked
}

func (dao *BreakDao) SelectLikedIdsBy(acquiree uint64, pagination model.Pagination) ([]uint64, error) {
	var ids []uint64
	sql := "select break_id from break_like where account_id=? order by create_time desc limit ?,?"
	err := dao.Tx.Raw(sql, acquiree, (pagination.Page-1)*pagination.Size, pagination.Size).Scan(&ids).Error
	return ids, err
}
