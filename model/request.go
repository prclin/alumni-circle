package model

// Pagination 分页
type Pagination struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1"`
}
