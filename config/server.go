package config

type Server struct {
	Port        int
	ContextPath string
	//服务器最大连接数
	MaxConnections uint
	//最大长连接数
	MaxKeepAliveRequests uint
}

var DefaultServer = &Server{
	Port:                 8080,
	ContextPath:          "/",
	MaxConnections:       8192,
	MaxKeepAliveRequests: 1000,
}
