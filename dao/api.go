package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
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

func (ad *APIDao) UpdateBy(tapi model.TAPI) error {
	sql := "update api set name=?,method=?,path=?,description=?,state=?,extra=? where id=?"
	return ad.Tx.Exec(sql, tapi.Name, tapi.Method, tapi.Path, tapi.Description, tapi.State, tapi.Extra, tapi.Id).Error
}

func (ad *APIDao) DeleteBindingByAPIId(apiId uint32) error {
	sql := "delete from api_binding where api_id=?"
	return ad.Tx.Exec(sql, apiId).Error
}

func (ad *APIDao) DeleteById(id uint32) error {
	sql := "delete from api where id=?"
	return ad.Tx.Exec(sql, id).Error
}

func (ad *APIDao) SelectPageBy(offset, size int) ([]model.TAPI, error) {
	var apis []model.TAPI
	sql := "select id, name, method, path, description, state, extra, create_time, update_time from api limit ?,?"
	err := ad.Tx.Raw(sql, offset, size).Scan(&apis).Error
	if apis == nil {
		apis = make([]model.TAPI, 0, 0)
	}
	return apis, err
}

func (ad *APIDao) BatchInsertBindingBy(bindings []model.TAPIBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into api_binding(role_id, api_id) values" //goland报错，忽略
	params := make([]interface{}, 0, len(bindings)*2)
	for _, binding := range bindings {
		sql += "(?,?),"
		params = append(params, binding.RoleId, binding.APIId)
	}
	sql = strings.TrimSuffix(sql, ",")
	return ad.Tx.Exec(sql, params...).Error
}

func (ad *APIDao) SelectCountByIds(ids []uint32) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var count int
	sql := "select count(id) from api where id in ?"
	err := ad.Tx.Raw(sql, ids).First(&count).Error
	return count, err
}
