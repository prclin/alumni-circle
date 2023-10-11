package controller

import (
	"fmt"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"math/rand"
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
	code := rand.Intn(9999)
	request.SetTemplateParam(fmt.Sprintf("{\"code\":%v}", code))
	request.SetTemplateCode("SMS_154950909")
	global.SMSClient.SendSms(request)
	model.Ok(context, struct{}{})
}
