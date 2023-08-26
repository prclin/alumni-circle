package model

import "time"

// TImage 图片表
type TImage struct {
	Id         uint64
	URL        string
	Extra      *string
	CreateTime time.Time
	UpdateTime time.Time
}

// Photo 照片
type Photo struct {
	URL   string `json:"url"`
	Order uint8  `json:"order"`
}

// TPhotoBinding 照片绑定
type TPhotoBinding struct {
	AccountId uint64 `json:"-"`
	ImageId   uint64 `json:"image_id" binding:"required"`
	Order     uint8  `json:"order" binding:"required"`
}
