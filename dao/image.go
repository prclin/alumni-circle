package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
)

type PhotoDao struct {
	Tx *gorm.DB
}

func NewPhotoDao(tx *gorm.DB) *PhotoDao {
	return &PhotoDao{Tx: tx}
}

func (pd *PhotoDao) SelectPhotosByAccountId(accountId uint64) ([]model.Photo, error) {
	var photos []model.Photo
	sql := "select i.url, pb.`order` from photo_binding as pb left join image as i on pb.image_id=i.id where pb.account_id=? order by pb.`order`"
	err := pd.Tx.Raw(sql, accountId).Scan(&photos).Error
	return photos, err
}

func (pd *PhotoDao) DeleteByAccountId(accountId uint64) error {
	sql := "delete from photo_binding where account_id=?"
	return pd.Tx.Exec(sql, accountId).Error
}

func (pd *PhotoDao) BatchInsertBy(bindings []model.TPhotoBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into photo_binding(account_id, image_id, `order`) values" //此处为goland报错
	params := make([]interface{}, 0, len(bindings))
	for _, binding := range bindings {
		sql += "(?,?,?),"
		params = append(params, binding.AccountId, binding.ImageId, binding.Order)
	}
	sql = strings.TrimSuffix(sql, ",")
	return pd.Tx.Exec(sql, params...).Error
}
