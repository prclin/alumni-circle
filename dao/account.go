package dao

import (
	. "github.com/prclin/alumni-circle/model/po"
	"gorm.io/gorm"
)

type AccountDao struct {
	Tx *gorm.DB
}

func NewAccountDao(tx *gorm.DB) *AccountDao {
	return &AccountDao{Tx: tx}
}

func (ad *AccountDao) InsertByAccount(account TAccount) (uint64, error) {
	var id uint64
	sql := "insert into account(email, password) value (?,?)"
	//插入用户
	if err := ad.Tx.Exec(sql, account.Email, account.Password).Error; err != nil {
		return 0, err
	}
	//查询主键
	if err := ad.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (ad *AccountDao) SelectByEmail(email string) (TAccount, error) {
	var account TAccount
	sql := "select id, phone, email, password, state, extra, create_time, update_time from account where email=?"
	err := ad.Tx.Raw(sql, email).First(&account).Error
	return account, err
}

type AccountInfoDao struct {
	Tx *gorm.DB
}

func NewAccountInfoDao(tx *gorm.DB) *AccountInfoDao {
	return &AccountInfoDao{Tx: tx}
}

func (aid *AccountInfoDao) InsertByAccountInfo(accountInfo TAccountInfo) error {
	sql := "insert into account_info(id, avatar_url, nickname) value (?,?,?)"
	return aid.Tx.Exec(sql, accountInfo.Id, accountInfo.AvatarURL, accountInfo.Nickname).Error
}

func (aid *AccountInfoDao) SelectById(id uint64) (TAccountInfo, error) {
	var ai TAccountInfo
	sql := "select id, campus_id, avatar_url, nickname, sex, birthday, follow_count, follower_count, extra, create_time, update_time from account_info where id=?"
	err := aid.Tx.Raw(sql, id).First(&ai).Error
	return ai, err
}

func (aid *AccountInfoDao) UpdateBy(info TAccountInfo) error {
	sql := "update account_info set campus_id=?,avatar_url=?,nickname=?,sex=?,birthday=?,extra=? where id=?"
	return aid.Tx.Exec(sql, info.CampusId, info.AvatarURL, info.Nickname, info.Sex, info.Birthday, info.Extra, info.Id).Error
}
