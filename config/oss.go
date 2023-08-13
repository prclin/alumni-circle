package config

// OSS OSS配置
type OSS struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyIdy"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
	Path            string `yaml:"path"`
	URL             string `yaml:"url"`
}

var DefaultOSS = &OSS{
	Endpoint:        "https://oss-cn-chengdu.aliyuncs.com",
	AccessKeyId:     "LTAI5tDU1MgTC49kLbyA9ZCL",
	AccessKeySecret: "WYLTdfKU39btlK7eqPiVgP7to3Nd9C",
	BucketName:      "alumni-circle",
	Path:            "development/",
	URL:             "https://alumni-circle.oss-cn-chengdu.aliyuncs.com",
}
