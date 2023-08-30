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
