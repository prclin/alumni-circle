package service

import (
	"github.com/prclin/alumni-circle/dao"
	. "github.com/prclin/alumni-circle/global"
	"gopkg.in/gomail.v2"
	"math/rand"
	"strconv"
	"time"
)

var dialer = gomail.NewDialer(Configuration.Email.Host, Configuration.Email.Port, Configuration.Email.Username, Configuration.Email.Password)

// SendVerifyEmail 发送校验邮件
func SendVerifyEmail(to string) error {
	//生成验证码
	captcha := rand.Intn(999999)
	//存储校验码到redis
	err := dao.SetString("captcha:"+to, strconv.Itoa(captcha), time.Minute*5)
	if err != nil {
		Logger.Warn(err)
		return err
	}
	//发送邮件
	message := gomail.NewMessage()
	message.SetHeader("From", Configuration.Email.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Verify Email")
	message.SetBody("text/html", "验证码:"+strconv.Itoa(captcha))
	err = dialer.DialAndSend(message)
	if err != nil {
		Logger.Warn(err)
		//删除captcha
		dao.DeleteKey("captcha:" + to)
		return err
	}
	return nil
}
