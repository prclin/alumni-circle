package controller

import (
	"fmt"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"math/rand"
	"strconv"
	"time"
)

// 注册路由
func init() {
	sms := core.ContextRouter.Group("/sms")
	sms.GET("/:phone", GetCaptcha)
}

// GetCaptcha 获取账户登录验证码
func GetCaptcha(context *gin.Context) {
	//获取电话号码
	phone := context.Param("phone")

	request := &dysmsapi.SendSmsRequest{}
	request.SetPhoneNumbers(phone)
	request.SetSignName("阿里云短信测试")
	code := rand.Intn(999999)
	err := dao.SetString("captcha:"+phone, strconv.Itoa(code), 60*time.Second)
	if err != nil {
		model.Client(context)
		return
	}
	request.SetTemplateParam(fmt.Sprintf("{\"code\":%v}", code))
	request.SetTemplateCode("SMS_154950909")
	_, err = global.SMSClient.SendSms(request)
	if err != nil {
		global.Logger.Debug(err)
		model.Server(context)
		return
	}
	model.Ok(context, struct{}{})
}
