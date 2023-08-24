package po

type TPhotoBinding struct {
	AccountId uint64 `json:"-"`
	ImageId   uint64 `json:"image_id" binding:"required"`
	Order     uint8  `json:"order" binding:"required"`
}
