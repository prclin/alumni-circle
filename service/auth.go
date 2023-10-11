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

func PhoneCaptchaSignIn(phone, captcha string) model.Response[*string] {
	//查账户
	ad := dao.NewAccountDao(Datasource)
	tAccount, err := ad.SelectByPhone(phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}

	//账户不存在
	if err == gorm.ErrRecordNotFound {
		resp := PhoneCaptchaSignUp(phone, captcha)
		if resp.Code != http.StatusOK {
			return resp
		}
		//再查账户
		tAccount, err = ad.SelectByPhone(phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			Logger.Debug(err)
			return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
		}
	} else { //存在
		//获取验证码
		resCap, err := dao.GetString("captcha:" + phone)
		if err != nil && !errors.Is(err, redis.Nil) {
			Logger.Debug(err)
			return model.Response[*string]{Code: 500, Message: "服务器内部错误"}
		}

		//校验
		if resCap == "" || captcha != resCap {
			Logger.Info("验证码错误:", "client", captcha, " ", "server-", resCap)
			return model.Response[*string]{Code: http.StatusBadRequest, Message: "验证码错误"}
		}
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

	return model.Response[*string]{Code: http.StatusOK, Message: "登录成功", Data: &token}
}

func PhoneCaptchaSignUp(phone, captcha string) model.Response[*string] {
	//获取验证码
	resCap, err := dao.GetString("captcha:" + phone)
	if err != nil && !errors.Is(err, redis.Nil) {
		Logger.Debug(err)
		return model.Response[*string]{Code: 500, Message: "服务器内部错误"}
	}

	//校验
	if resCap == "" || captcha != resCap {
		Logger.Info("验证码错误:", "client", captcha, " ", "server-", resCap)
		return model.Response[*string]{Code: http.StatusBadRequest, Message: "验证码错误"}
	}

	//插入用户
	tx := Datasource.Begin() //开启事务
	defer tx.Commit()
	accountDao := dao.NewAccountDao(tx)
	id, err := accountDao.InsertByAccount(model.TAccount{Phone: phone})
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}

	//初始化账户信息
	info := model.TAccountInfo{Id: id, Nickname: fmt.Sprintf("用户-%v", id), AvatarURL: "默认"}
	accountInfoDao := dao.NewAccountInfoDao(tx)
	err = accountInfoDao.InsertByAccountInfo(info)
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}
	return model.Response[*string]{Code: http.StatusOK, Message: "注册成功"}
}

func PhonePasswordSignIn(phone, password string) model.Response[*string] {
	//查账户
	ad := dao.NewAccountDao(Datasource)
	tAccount, err := ad.SelectByPhone(phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}
	//账户不存在
	if err == gorm.ErrRecordNotFound {
		Logger.Debug(err)
		return model.Response[*string]{Code: http.StatusBadRequest, Message: "账户不存在!"}
	}
	//比对密码
	if util.MD5([]byte(password)) != tAccount.Password {
		return model.Response[*string]{Code: http.StatusBadRequest, Message: "密码错误!"}
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

func RevokeRoleAllocation(accountIds []uint64, roleIds []uint32) error {
	//映射binding
	bindings := make([]model.TRoleBinding, 0, len(accountIds)*len(roleIds))
	for _, accountId := range accountIds {
		for _, roleId := range roleIds {
			bindings = append(bindings, model.TRoleBinding{AccountId: accountId, RoleId: roleId})
		}
	}
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	roleDao := dao.NewRoleDao(tx)
	//删除
	err := roleDao.DeleteBindingBy(bindings)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func AllocateRole(accountIds []uint64, roleIds []uint32) error {
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	//account ids 是否正确
	accountDao := dao.NewAccountDao(tx)
	accountCount, err := accountDao.SelectCountByIds(accountIds)
	if err != nil || accountCount != len(accountIds) {
		tx.Rollback()
		return util.Ternary(err != nil, err, errors.New("部分账户不存在"))
	}
	//role ids是否正确
	roleDao := dao.NewRoleDao(tx)
	roleCount, err := roleDao.SelectCountByIds(roleIds)
	if err != nil || roleCount != len(roleIds) {
		tx.Rollback()
		return util.Ternary(err != nil, err, errors.New("部分角色不存在"))
	}
	//映射binding
	bindings := make([]model.TRoleBinding, 0, len(accountIds)*len(roleIds))
	for _, accountId := range accountIds {
		for _, roleId := range roleIds {
			bindings = append(bindings, model.TRoleBinding{AccountId: accountId, RoleId: roleId})
		}
	}
	//分配
	err = roleDao.BatchInsertBindingBy(bindings)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func RevokeAPIAllocation(roleIds, apiIds []uint32) error {
	//映射binding
	bindings := make([]model.TAPIBinding, 0, len(roleIds)*len(apiIds))
	for _, roleId := range roleIds {
		for _, apiId := range apiIds {
			bindings = append(bindings, model.TAPIBinding{RoleId: roleId, APIId: apiId})
		}
	}
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	apiDao := dao.NewAPIDao(tx)
	err := apiDao.DeleteBindingBy(bindings)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func AllocateAPI(roleIds, apiIds []uint32) error {
	//事务
	tx := Datasource.Begin()
	defer tx.Commit()
	//role ids 是否正确
	roleDao := dao.NewRoleDao(tx)
	roleCount, err := roleDao.SelectCountByIds(roleIds)
	if err != nil || roleCount != len(roleIds) {
		tx.Rollback()
		return util.Ternary(err != nil, err, errors.New("部分角色不存在"))
	}
	//api ids是否正确
	apiDao := dao.NewAPIDao(tx)
	apiCount, err := apiDao.SelectCountByIds(apiIds)
	if err != nil || apiCount != len(apiIds) {
		tx.Rollback()
		return util.Ternary(err != nil, err, errors.New("部分接口不存在"))
	}
	//映射binding
	bindings := make([]model.TAPIBinding, 0, len(roleIds)*len(apiIds))
	for _, roleId := range roleIds {
		for _, apiId := range apiIds {
			bindings = append(bindings, model.TAPIBinding{RoleId: roleId, APIId: apiId})
		}
	}
	//批量插入
	err = apiDao.BatchInsertBindingBy(bindings)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetRoleList(pagination model.Pagination) ([]model.TRole, error) {
	roleDao := dao.NewRoleDao(Datasource)
	return roleDao.SelectPageBy((pagination.Page-1)*pagination.Size, pagination.Size)
}

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
	if resCap == "" || captcha != resCap {
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
