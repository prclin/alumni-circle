package config

type SMS struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
}

var DefaultSMS = &SMS{Endpoint: "dysmsapi.aliyuncs.com", AccessKeyId: "", AccessKeySecret: ""}
