package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
)

type CampusDao struct {
	Tx *gorm.DB
}

func NewCampusDao(tx *gorm.DB) *CampusDao {
	return &CampusDao{Tx: tx}
}

func (dao CampusDao) SelectById(id uint32) (*model.TCampus, error) {
	var campus *model.TCampus
	sql := "select id,name,extra,create_time,update_time from campus where id = ?"
	err := dao.Tx.Exec(sql, id).First(campus).Error
	return campus, err
}
