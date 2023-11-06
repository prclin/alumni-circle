package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
)

type ImageDao struct {
	Tx *gorm.DB
}

func NewImageDao(tx *gorm.DB) *ImageDao {
	return &ImageDao{Tx: tx}
}

func (imageDao *ImageDao) InsertBy(image model.TImage) (uint64, error) {
	var id uint64
	sql := "insert into image(url, extra) value (?,?)"
	//插入
	if err := imageDao.Tx.Exec(sql, image.URL, image.Extra).Error; err != nil {
		return 0, err
	}
	//获取主键
	if err := imageDao.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (imageDao *ImageDao) SelectById(id uint64) (model.TImage, error) {
	var image model.TImage
	sql := "select id, url, extra, create_time, update_time from image where id=?"
	err := imageDao.Tx.Raw(sql, id).First(&image).Error
	return image, err
}

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

type ShotDao struct {
	Tx *gorm.DB
}

func NewShotDao(tx *gorm.DB) *ShotDao {
	return &ShotDao{Tx: tx}
}

func (sd *ShotDao) BatchInsertBy(bindings []model.TShotBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into shot_binding(break_id, image_id, `order`) values" //此处为goland报错
	params := make([]interface{}, 0, len(bindings)*3)
	for _, binding := range bindings {
		sql += "(?,?,?),"
		params = append(params, binding.BreakId, binding.ImageId, binding.Order)
	}
	sql = strings.TrimSuffix(sql, ",")
	return sd.Tx.Exec(sql, params...).Error
}

func (sd *ShotDao) SelectShotsByBreakId(breakId uint64) ([]model.Shot, error) {
	var shots []model.Shot
	sql := "select i.id, i.url, i.extra, i.create_time, i.update_time, sb.`order` from shot_binding as sb left join image as i on sb.image_id=i.id where sb.break_id=? order by sb.`order`"
	err := sd.Tx.Raw(sql, breakId).Scan(&shots).Error
	if shots == nil {
		shots = make([]model.Shot, 0, 0)
	}
	return shots, err
}
