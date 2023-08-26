package service

import (
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"mime/multipart"
	"path"
)

func UploadFileToOSS(bucketName string, dir string, fileHeader *multipart.FileHeader) (model.TImage, error) {
	//获取bucket
	bucket, err := global.OSSClient.Bucket(bucketName)
	if err != nil {
		return model.TImage{}, err
	}
	//重命名文件
	fileHeader.Filename = path.Join(dir, fileHeader.Filename)
	//打开文件
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return model.TImage{}, err
	}
	//上传文件
	err = bucket.PutObject(fileHeader.Filename, file)
	if err != nil {
		return model.TImage{}, err
	}
	//保存图片信息
	image := model.TImage{URL: "https://" + bucketName + "." + global.Configuration.OSS.EndPoint + "/" + fileHeader.Filename}
	tx := global.Datasource.Begin()
	defer tx.Commit()
	imageDao := dao.NewImageDao(tx)
	id, err := imageDao.InsertBy(image) //插入记录
	if err != nil {
		tx.Rollback()
		return model.TImage{}, err
	}
	//获取新插入照片信息
	tImage, err := imageDao.SelectById(id)
	if err != nil {
		tx.Rollback()
		return model.TImage{}, err
	}
	return tImage, err

}
