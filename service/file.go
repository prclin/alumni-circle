package service

import (
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	"github.com/prclin/alumni-circle/util"
	"mime/multipart"
	"path"
	"time"
)

func ImageUpload(filename string, file multipart.File) (image *entity.Image, err error) {
	// 拼接文件名,格式为 development/yyyy/MM/dd/hh:mm:ss-fileNameMD5.ext
	filename = global.Configuration.OSS.Path + time.Now().Format("2006/01/02/15:04:05") + "-" + util.StringMD5(filename) + path.Ext(filename)
	// 上传文件
	err = global.OSSBucket.PutObject(filename, file)
	if err != nil {
		return nil, err
	}
	// 创建image对象存入数据库
	image = new(entity.Image)
	image.Url = global.Configuration.OSS.URL + "/" + filename
	err = entity.CreateImage(image)
	if err != nil {
		return nil, err
	}
	return image, err
}

func ImageExist(id int) bool {
	image := &entity.Image{
		Id: id,
	}
	if err := entity.GetImage(image); err != nil {
		return false
	}
	return true
}
