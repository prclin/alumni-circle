package service

import (
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

func GetAccountInfo(acquirer uint64, acquiree uint64) (model.Account, error) {
	//获取账户信息
	aid := dao.NewAccountInfoDao(Datasource)
	//获取账户信息
	info, err := aid.SelectById(acquiree)
	if err != nil {
		return model.Account{}, err
	}
	//获取关系
	fd := dao.NewFollowDao(Datasource)
	followed, err := fd.IsFollowed(acquirer, acquirer)
	if err != nil {
		return model.Account{}, err
	}
	return model.Account{Info: info, IsFollowed: followed}, nil
}

func UpdateAccountInfo(info model.TAccountInfo) error {
	aid := dao.NewAccountInfoDao(Datasource)
	return aid.UpdateBy(info)
}

func GetPhotoWall(accountId uint64) ([]model.Photo, error) {
	pd := dao.NewPhotoDao(Datasource)
	return pd.SelectPhotosByAccountId(accountId)
}

func UpdateAccountPhoto(accountId uint64, bindings []model.TPhotoBinding) error {
	//开启事务
	tx := Datasource.Begin()
	defer tx.Commit()
	pd := dao.NewPhotoDao(tx)
	//删除原先的照片
	err := pd.DeleteByAccountId(accountId)
	if err != nil {
		tx.Rollback()
		return err
	}
	//插入新照片
	for i := 0; i < len(bindings); i++ {
		bindings[i].AccountId = accountId
	}
	err = pd.BatchInsertBy(bindings)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
