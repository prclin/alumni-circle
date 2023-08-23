package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	. "github.com/prclin/alumni-circle/model/po"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

func EmailSignUp(account TAccount, captcha string) Response[any] {
	//获取验证码
	resCap, err := dao.GetString(fmt.Sprintf("captcha:%v", account.Email))
	if err != nil && !errors.Is(err, redis.Nil) {
		Logger.Debug(err)
		return Response[any]{Code: 500, Message: "服务器内部错误"}
	}

	//校验
	if captcha != resCap {
		Logger.Info("验证码错误:", "client", captcha, " ", "server-", resCap)
		return Response[any]{Code: http.StatusBadRequest, Message: "验证码错误"}
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
		return Response[any]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}

	//初始化账户信息
	info := TAccountInfo{Id: id, Nickname: fmt.Sprintf("用户-%v", id), AvatarURL: "默认"}
	accountInfoDao := dao.NewAccountInfoDao(tx)
	err = accountInfoDao.InsertByAccountInfo(info)
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return Response[any]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}
	return Response[any]{Code: http.StatusOK, Message: "注册成功"}
}

func EmailSignIn(email, password string) Response[*string] {
	//查密码
	ad := dao.NewAccountDao(Datasource)
	tAccount, err := ad.SelectByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}

	//比对密码
	if util.MD5([]byte(password)) != tAccount.Password {
		return Response[*string]{Code: http.StatusBadRequest, Message: "邮箱或密码错误!"}
	}

	//获取账户角色
	rd := dao.NewRoleDao(Datasource)
	bindings, err := rd.SelectBindingByAccountId(tAccount.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Debug(err)
		return Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}

	//映射
	roleIds := make([]uint32, 0, len(bindings))
	for _, binding := range bindings {
		roleIds = append(roleIds, binding.RoleId)
	}

	//生成token
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Id:               tAccount.Id,
		RoleIds:          roleIds,
	}
	token, err := util.GenerateToken(claims)
	if err != nil {
		Logger.Debug(err)
		return Response[*string]{Code: http.StatusInternalServerError, Message: "服务器内部错误!"}
	}
	return Response[*string]{Code: http.StatusOK, Message: "登录成功！", Data: &token}
}
