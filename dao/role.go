package dao

import (
	. "github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
	"strings"
)

type RoleDao struct {
	Tx *gorm.DB
}

func NewRoleDao(tx *gorm.DB) *RoleDao {
	return &RoleDao{Tx: tx}
}

func (rd *RoleDao) SelectBindingByAccountId(accountId uint64) ([]TRoleBinding, error) {
	var rb []TRoleBinding
	sql := "select account_id, role_id, extra, create_time, update_time from role_binding where account_id=?"
	err := rd.Tx.Raw(sql, accountId).Scan(&rb).Error
	return rb, err
}

func (rd *RoleDao) InsertBy(role TRole) (uint32, error) {
	var id uint32
	sql := "insert into role(name, identifier, description, state) value (?,?,?,?)"
	//插入
	if err := rd.Tx.Exec(sql, role.Name, role.Identifier, role.Description, role.State).Error; err != nil {
		return 0, err
	}
	//获取主键
	if err := rd.Tx.Raw("select LAST_INSERT_ID()").First(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (rd *RoleDao) SelectById(id uint32) (TRole, error) {
	var role TRole
	sql := "select id, name, identifier, description, state, create_time,update_time from role where id=?"
	err := rd.Tx.Raw(sql, id).First(&role).Error
	return role, err
}

func (rd *RoleDao) UpdateBy(tRole TRole) error {
	sql := "update role set name=?,identifier=?,description=?,state=? where id=?"
	return rd.Tx.Exec(sql, tRole.Name, tRole.Identifier, tRole.Description, tRole.State, tRole.Id).Error
}

func (rd *RoleDao) DeleteBindingByRoleId(roleId uint32) error {
	sql := "delete from role_binding where role_id=?"
	return rd.Tx.Exec(sql, roleId).Error
}

func (rd *RoleDao) DeleteById(id uint32) error {
	sql := "delete from role where id=?"
	return rd.Tx.Exec(sql, id).Error
}

func (rd *RoleDao) SelectPageBy(offset, size int) ([]TRole, error) {
	var roles []TRole
	sql := "select id, name, identifier, description, state, create_time, update_time from role limit ?,?"
	err := rd.Tx.Raw(sql, offset, size).Scan(&roles).Error
	if roles == nil {
		roles = make([]TRole, 0, 0)
	}
	return roles, err
}

func (rd *RoleDao) SelectCountByIds(ids []uint32) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var count int
	sql := "select count(id) from role where id in ?"
	err := rd.Tx.Raw(sql, ids).First(&count).Error
	return count, err
}

func (rd *RoleDao) BatchInsertBindingBy(bindings []TRoleBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "insert into role_binding(account_id, role_id) values "
	params := make([]interface{}, 0, len(bindings)*2)
	for _, binding := range bindings {
		sql += "(?,?),"
		params = append(params, binding.AccountId, binding.RoleId)
	}
	sql = strings.TrimSuffix(sql, ",")
	return rd.Tx.Exec(sql, params...).Error
}

func (rd *RoleDao) DeleteBindingBy(bindings []TRoleBinding) error {
	if len(bindings) == 0 {
		return nil
	}
	sql := "delete from role_binding where" //goland报错，忽略
	params := make([]interface{}, 0, len(bindings)*2)
	for _, binding := range bindings {
		sql += " (account_id=? and role_id=?) or"
		params = append(params, binding.AccountId, binding.RoleId)
	}
	sql = strings.TrimSuffix(sql, "or")
	return rd.Tx.Exec(sql, params...).Error
}
