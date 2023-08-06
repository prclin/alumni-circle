package config

/*
Zap 配置结构
*/
type Zap struct {
	//Level 日志输出级别，可选值见zapcore.Level
	Level string
}

var DefaultZap = &Zap{
	Level: "DEBUG",
}
