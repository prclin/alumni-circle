package service

import (
	"errors"
	"fmt"
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	. "github.com/prclin/alumni-circle/model/entity"
	. "github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/util"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func EmailSignUp(account Account, captcha string) Response[any] {
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
	info := AccountInfo{Id: id, Nickname: fmt.Sprintf("用户-%v", id), AvatarURL: "默认"}
	accountInfoDao := dao.NewAccountInfoDao(tx)
	err = accountInfoDao.InsertByAccountInfo(info)
	if err != nil {
		tx.Rollback()
		Logger.Debug(err)
		return Response[any]{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	}
	return Response[any]{Code: http.StatusOK, Message: "注册成功"}
}
