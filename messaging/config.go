package messaging

// Config message broker配置
type Config struct {
	//应用消息前缀
	AppDestinationPrefix string
	//代理消息前缀
	BrokerDestinationPrefix string
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	AppDestinationPrefix:    "/app",
	BrokerDestinationPrefix: "/topic",
}
