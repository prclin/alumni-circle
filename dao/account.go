package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
)

type AccountDao struct {
	Tx *gorm.DB
}

func NewAccountDao(tx *gorm.DB) *AccountDao {
	return &AccountDao{Tx: tx}
}

func (ad *AccountDao) InsertByAccount(account model.TAccount) (uint64, error) {
	var id uint64
	sql := "insert into account(phone,email,password,extra) value (?,?,?,?)"
	//插入用户
	if err := ad.Tx.Exec(sql, account.Phone, account.Email, account.Password, account.Extra).Error; err != nil {
		return 0, err
	}
	//查询主键
	if err := ad.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (ad *AccountDao) SelectByEmail(email string) (model.TAccount, error) {
	var account model.TAccount
	sql := "select id, phone, email, password, state, extra, create_time, update_time from account where email=?"
	err := ad.Tx.Raw(sql, email).First(&account).Error
	return account, err
}

func (ad *AccountDao) SelectCountByIds(ids []uint64) (int, error) {
	var count int
	sql := "select count(id) from account where id in ?"
	err := ad.Tx.Raw(sql, ids).First(&count).Error
	return count, err
}

func (ad *AccountDao) Exist(id uint64) (bool, error) {
	var exist bool
	sql := "select count(id) from account where id = ?"
	err := ad.Tx.Raw(sql, id).First(&exist).Error
	return exist, err
}

func (ad *AccountDao) SelectByPhone(phone string) (model.TAccount, error) {
	var account model.TAccount
	sql := "select id, phone, email, password, state, extra, create_time, update_time from account where phone=?"
	err := ad.Tx.Raw(sql, phone).First(&account).Error
	return account, err
}

type AccountInfoDao struct {
	Tx *gorm.DB
}

func NewAccountInfoDao(tx *gorm.DB) *AccountInfoDao {
	return &AccountInfoDao{Tx: tx}
}

func (aid *AccountInfoDao) InsertByAccountInfo(accountInfo model.TAccountInfo) error {
	sql := "insert into account_info(id, avatar_url, nickname) value (?,?,?)"
	return aid.Tx.Exec(sql, accountInfo.Id, accountInfo.AvatarURL, accountInfo.Nickname).Error
}

func (aid *AccountInfoDao) SelectById(id uint64) (*model.TAccountInfo, error) {
	var ai *model.TAccountInfo
	sql := "select id, campus_id, nickname, avatar_url,background_url, sex, birthday, brief, mbti_result_id, follow_count, follower_count, friend_count, extra, create_time, update_time from account_info where id=?"
	err := aid.Tx.Raw(sql, id).First(&ai).Error
	return ai, err
}

func (aid *AccountInfoDao) UpdateBy(info model.TAccountInfo) error {
	sql := "update account_info set nickname=?,avatar_url=?,background_url=?,sex=?,brief=?,birthday=?,extra=? where id=?"
	return aid.Tx.Exec(sql, info.Nickname, info.AvatarURL, info.BackgroundURL, info.Sex, info.Brief, info.Birthday, info.Extra, info.Id).Error
}

func (aid *AccountInfoDao) UpdateMBTIResultIdById(id uint64, mbtiResultId string) error {
	sql := "update account_info set mbti_result_id=? where id=?"
	return aid.Tx.Exec(sql, mbtiResultId, id).Error
}

type FollowDao struct {
	Tx *gorm.DB
}

func NewFollowDao(tx *gorm.DB) *FollowDao {
	return &FollowDao{Tx: tx}
}

func (fd *FollowDao) IsFollowed(follower, followee uint64) (bool, error) {
	var followed bool
	sql := "select count(*) from follow where follower_id=? and followee_id=?"
	err := fd.Tx.Raw(sql, follower, followee).First(&followed).Error
	return followed, err
}

func (fd *FollowDao) InsertBy(follow model.TFollow) error {
	sql := "insert into follow(follower_id, followee_id, extra) value (?,?,?)"
	return fd.Tx.Exec(sql, follow.FollowerId, follow.FolloweeId, follow.Extra).Error
}

func (fd *FollowDao) DeleteBy(follow model.TFollow) error {
	sql := "delete from follow where follower_id=? and followee_id=?"
	return fd.Tx.Exec(sql, follow.FollowerId, follow.FolloweeId).Error
}
