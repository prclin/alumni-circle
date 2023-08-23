package service

import (
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
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
