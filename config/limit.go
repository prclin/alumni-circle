package config

// Limit Limit配置，包含对参数数量,大小等限制
type Limit struct {
	// 大小限制
	Size *Size
}
type Size struct {
	// 一篇课间包含图片数量限制
	PictureInBreak int
}

var DefaultLimit = &Limit{
	Size: &Size{
		PictureInBreak: 9,
	},
}
