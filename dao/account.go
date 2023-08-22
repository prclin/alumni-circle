package dao

import (
	. "github.com/prclin/alumni-circle/model/entity"
	"gorm.io/gorm"
)

type AccountDao struct {
	Tx *gorm.DB
}

func NewAccountDao(tx *gorm.DB) AccountDao {
	return AccountDao{Tx: tx}
}

func (ad *AccountDao) InsertByAccount(account Account) (uint64, error) {
	var id uint64
	sql := "insert into account(email, password) value (?,?)"
	//插入用户
	if err := ad.Tx.Exec(sql, account.Email, account.Password).Error; err != nil {
		return 0, err
	}
	//查询主键
	if err := ad.Tx.Raw("select LAST_INSERT_ID()").Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

type AccountInfoDao struct {
	Tx *gorm.DB
}

func NewAccountInfoDao(tx *gorm.DB) *AccountInfoDao {
	return &AccountInfoDao{Tx: tx}
}

func (aid *AccountInfoDao) InsertByAccountInfo(accountInfo AccountInfo) error {
	sql := "insert into account_info(id, avatar_url, nickname) value (?,?,?)"
	return aid.Tx.Exec(sql, accountInfo.Id, accountInfo.AvatarURL, accountInfo.Nickname).Error
}
