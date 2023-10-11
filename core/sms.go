package core

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/prclin/alumni-circle/global"
)

// initSMS 初始化阿里云sms客户端
func initSMS() {
	smsConfig := global.Configuration.SMS
	config := &openapi.Config{ // 您的AccessKey ID
		AccessKeyId: &smsConfig.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &smsConfig.AccessKeySecret,
	}
	config.Endpoint = tea.String(smsConfig.Endpoint)
	client, err := dysmsapi.NewClient(config)
	if err != nil {
		global.Logger.Fatal(err)
	}
	global.SMSClient = client
}
