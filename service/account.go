package service

import (
	"errors"
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/util"
)

func RevokeFollow(follow model.TFollow) error {
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	//取关
	followDao := dao.NewFollowDao(tx)
	err := followDao.DeleteBy(follow)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func FollowAccount(follow model.TFollow) error {
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	accountDao := dao.NewAccountDao(tx)
	//验证被关注账户
	exist, err := accountDao.Exist(follow.FolloweeId)
	if err != nil || !exist {
		tx.Rollback()
		return util.Ternary(err != nil, err, errors.New("被关注账户不存在"))
	}
	//关注
	followDao := dao.NewFollowDao(tx)
	err = followDao.InsertBy(follow)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetAccountTag(id uint64) ([]model.TTag, error) {
	td := dao.NewTagDao(Datasource)
	tags, err := td.SelectEnabledByAccountId(id)
	if tags == nil {
		tags = make([]model.TTag, 0, 0)
	}
	return tags, err
}

func UpdateAccountTag(id uint64, tagIds []uint32) ([]model.TTag, error) {
	tx := Datasource.Begin()
	defer tx.Commit()
	td := dao.NewTagDao(tx)
	//删除原标签
	err := td.DeleteBindingByAccountId(id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//绑定新标签
	bindings := make([]model.TTagBinding, 0, len(tagIds))
	for _, tagId := range tagIds {
		bindings = append(bindings, model.TTagBinding{AccountId: id, TagId: tagId})
	}
	err = td.BatchInsertBindingBy(bindings)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//获取绑定后标签列表
	tags, err := td.SelectEnabledByIds(tagIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if tags == nil {
		tags = make([]model.TTag, 0, 0)
	}
	return tags, nil
}

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
