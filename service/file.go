package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	"github.com/prclin/alumni-circle/util"
	"mime/multipart"
	"path"
	"time"
)

// 初始全局化变量
const (
	// 配置阿里云OSS连接
	OSSEndpoint        = "https://oss-cn-chengdu-internal.aliyuncs.com"
	OSSAccessKeyId     = "LTAI5tDU1MgTC49kLbyA9ZCL"
	OSSAccessKeySecret = "WYLTdfKU39btlK7eqPiVgP7to3Nd9C"
	OSSBucketName      = "alumni-circle"
	OSSPath            = "development/"
	OSSUrl             = "https://alumni-circle.oss-cn-chengdu-internal.aliyuncs.com"
)

// OSS存储空间
var OSSBucket *oss.Bucket

// 初始化OSS连接
func init() {
	// 创建OSSClient实例
	client, err := oss.New(OSSEndpoint, OSSAccessKeyId, OSSAccessKeySecret)
	if err != nil {
		global.Logger.Fatal(err)
	}
	// 获取存储空间
	OSSBucket, err = client.Bucket(OSSBucketName)
	if err != nil {
		global.Logger.Fatal(err)
	}
}

func ImageUpload(filename string, file multipart.File) (image *entity.Image, err error) {
	// 拼接文件名,格式为 development/yyyy/MM/dd/hh:mm:ss-fileNameMD5.ext
	filename = OSSPath + time.Now().Format("2006/01/02/15:04:05") + "-" + util.StringMD5(filename) + path.Ext(filename)
	// 上传文件
	err = OSSBucket.PutObject(filename, file)
	if err != nil {
		return nil, err
	}
	// 创建image对象存入数据库
	image = new(entity.Image)
	image.Url = OSSUrl + "/" + filename
	err = entity.CreateImage(image)
	if err != nil {
		return nil, err
	}
	return image, err
}
