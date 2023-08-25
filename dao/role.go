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
