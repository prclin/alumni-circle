package config

var DefaultMBTI = &MBTI{
	API:   MBTIBaseAPI{},
	Token: "",
}

// MBTI mbti测试配置
type MBTI struct {
	API   MBTIBaseAPI
	Token string
}

// MBTIBaseAPI 基础接口配置
type MBTIBaseAPI struct {
	Sheet  string
	Submit string
}
