package dao

import (
	"github.com/prclin/alumni-circle/model/entity"
	"gorm.io/gorm"
)

type PhotoDao struct {
	Tx *gorm.DB
}

func NewPhotoDao(tx *gorm.DB) *PhotoDao {
	return &PhotoDao{Tx: tx}
}

func (pd *PhotoDao) SelectPhotosByAccountId(accountId uint64) ([]entity.Photo, error) {
	var photos []entity.Photo
	sql := "select i.url, pb.`order` from photo_binding as pb left join image as i on pb.image_id=i.id where pb.account_id=? order by pb.`order`"
	err := pd.Tx.Raw(sql, accountId).Scan(&photos).Error
	return photos, err
}
