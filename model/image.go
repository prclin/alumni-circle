package model

import "time"

// TImage 图片表
type TImage struct {
	Id         uint64    `json:"id"`
	URL        string    `json:"url"`
	Extra      *string   `json:"extra"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
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

// Shot 镜头
type Shot struct {
	TImage
	Order uint8 `json:"order"`
}

// TShotBinding 照片绑定
type TShotBinding struct {
	BreakId uint64
	ImageId uint64
	Order   uint8
}
