package service

import (
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	. "github.com/prclin/alumni-circle/model/po"
)

func GetAccountInfo(id uint64) (TAccountInfo, error) {
	//获取账户信息
	aid := dao.NewAccountInfoDao(Datasource)
	return aid.SelectById(id)
}

func UpdateAccountInfo(info TAccountInfo) error {
	aid := dao.NewAccountInfoDao(Datasource)
	return aid.UpdateBy(info)
}

func GetPhotoWall(accountId uint64) ([]entity.Photo, error) {
	pd := dao.NewPhotoDao(Datasource)
	return pd.SelectPhotosByAccountId(accountId)
}

func UpdateAccountPhoto(accountId uint64, bindings []TPhotoBinding) error {
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
