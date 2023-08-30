package dao

import (
	. "github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
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
