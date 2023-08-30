package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

func DeleteRole(id uint32) error {
	tx := Datasource.Begin()
	defer tx.Commit()
	roleDao := dao.NewRoleDao(tx)
	//删除绑定
	err := roleDao.DeleteBindingByRoleId(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	//删除角色
	err = roleDao.DeleteById(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdateRole(tRole model.TRole) (model.TRole, error) {
	tx := Datasource.Begin()
	defer tx.Commit()
	roleDao := dao.NewRoleDao(tx)
	//更新
	err := roleDao.UpdateBy(tRole)
	if err != nil {
		tx.Rollback()
		return model.TRole{}, err
	}
	//获取
	role, err := roleDao.SelectById(tRole.Id)
	if err != nil {
		tx.Rollback()
		return model.TRole{}, err
	}
	return role, nil
}

func CreateRole(tRole model.TRole) (model.TRole, error) {
	tx := Datasource.Begin()
	defer tx.Commit()
	roleDao := dao.NewRoleDao(tx)
	//插入
	id, err := roleDao.InsertBy(tRole)
	if err != nil {
		tx.Rollback()
		return model.TRole{}, err
	}
	//获取
	role, err := roleDao.SelectById(id)
	if err != nil {
		tx.Rollback()
		return model.TRole{}, err
	}
	return role, nil
}

func GetAPIList(pagination model.Pagination) ([]model.TAPI, error) {
	apiDao := dao.NewAPIDao(Datasource)
	return apiDao.SelectPageBy((pagination.Page-1)*pagination.Size, pagination.Size)
}

func DeleteAPI(id uint32) error {
	tx := Datasource.Begin()
	defer tx.Commit()
	apiDao := dao.NewAPIDao(tx)
	//解除绑定
	err := apiDao.DeleteBindingByAPIId(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	//删除api
	err = apiDao.DeleteById(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdateAPI(tapi model.TAPI) (model.TAPI, error) {
	tx := Datasource.Begin()
	defer tx.Commit()
	apiDao := dao.NewAPIDao(tx)
	//更新
	err := apiDao.UpdateBy(tapi)
	if err != nil {
		tx.Rollback()
		return model.TAPI{}, err
	}
	//获取最新的api
	api, err := apiDao.SelectById(tapi.Id)
	if err != nil {
		tx.Rollback()
		return model.TAPI{}, err
	}
	return api, nil
}

func CreateAPI(tapi model.TAPI) (model.TAPI, error) {
	tx := Datasource.Begin()
	defer tx.Commit()
	apiDao := dao.NewAPIDao(tx)
	//插入
	id, err := apiDao.InsertBy(tapi)
	if err != nil {
		tx.Rollback()
		return model.TAPI{}, err
	}
	//获取
	api, err := apiDao.SelectById(id)
	if err != nil {
		tx.Rollback()
		return model.TAPI{}, err
	}
	return api, nil
}

func EmailSignUp(account model.TAccount, captcha string) model.Response[any] {
	//获取验证码
	resCap, err := dao.GetString(fmt.Sprintf("captcha:%v", account.Email))
	if err != nil && !errors.Is(err, redis.Nil) {
		Logger.Debug(err)
		return model.Response[any]{Code: 500, Message: "服务器内部错误"}
	}

	//校验
	if captcha != resCap {
		Logger.Info("验证码错误:", "client", captcha, " ", "server-", resCap)
		return model.Response[any]{Code: http.StatusBadRequest, Message: "验证码错误"}
	}

	//密码加密
	account.Password = util.MD5([]byte(account.Password))

	//插入用户
	tx := Datasource.Begin() //开启事务
	defer tx.Commit()
	accountDao := dao.NewAccountDao(tx)
	id, err := accountDao.InsertByAccount(account)
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return model.Response[any]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}

	//初始化账户信息
	info := model.TAccountInfo{Id: id, Nickname: fmt.Sprintf("用户-%v", id), AvatarURL: "默认"}
	accountInfoDao := dao.NewAccountInfoDao(tx)
	err = accountInfoDao.InsertByAccountInfo(info)
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return model.Response[any]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}
	return model.Response[any]{Code: http.StatusOK, Message: "注册成功"}
}

func EmailSignIn(email, password string) model.Response[*string] {
	//查密码
	ad := dao.NewAccountDao(Datasource)
	tAccount, err := ad.SelectByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}

	//比对密码
	if util.MD5([]byte(password)) != tAccount.Password {
		return model.Response[*string]{Code: http.StatusBadRequest, Message: "邮箱或密码错误!"}
	}

	//获取账户角色
	rd := dao.NewRoleDao(Datasource)
	bindings, err := rd.SelectBindingByAccountId(tAccount.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}

	//映射
	roleIds := make([]uint32, 0, len(bindings))
	for _, binding := range bindings {
		roleIds = append(roleIds, binding.RoleId)
	}

	//生成token
	claims := model.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Id:               tAccount.Id,
		RoleIds:          roleIds,
	}
	token, err := util.GenerateToken(claims)
	if err != nil {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}
	return model.Response[*string]{Code: http.StatusOK, Message: "登录成功！", Data: &token}
}
