package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
)

type APIDao struct {
	Tx *gorm.DB
}

func NewAPIDao(tx *gorm.DB) *APIDao {
	return &APIDao{Tx: tx}
}

func (ad *APIDao) InsertBy(tapi model.TAPI) (uint32, error) {
	var id uint32
	sql := "insert into api(name, method, path, description, state, extra) VALUE (?,?,?,?,?,?)"
	if err := ad.Tx.Exec(sql, tapi.Name, tapi.Method, tapi.Path, tapi.Description, tapi.State, tapi.Extra).Error; err != nil {
		return 0, err
	}
	if err := ad.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (ad *APIDao) SelectById(id uint32) (model.TAPI, error) {
	var api model.TAPI
	sql := "select id, name, method, path, description, state, extra, create_time, update_time from api where id=?"
	err := ad.Tx.Raw(sql, id).First(&api).Error
	return api, err
}
