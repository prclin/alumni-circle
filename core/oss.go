package core

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/prclin/alumni-circle/global"
)

// initOSS 初始化OSS连接
func initOSS() {
	// 读取OSS配置
	OSSConfig := global.Configuration.OSS
	// 创建OSSClient实例
	client, err := oss.New(OSSConfig.Endpoint, OSSConfig.AccessKeyId, OSSConfig.AccessKeySecret)
	if err != nil {
		global.Logger.Fatal(err)
	}
	// 获取存储空间
	OSSBucket, err := client.Bucket(OSSConfig.BucketName)
	if err != nil {
		global.Logger.Fatal(err)
	}
	global.OSSBucket = OSSBucket
}
