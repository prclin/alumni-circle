package config

/*
Configuration 对应整个application.yaml配置文件
*/
type Configuration struct {
	Server     *Server
	Zap        *Zap
	Datasource *Datasource
	Email      *Email
	Redis      *Redis
}

var DefaultConfiguration = &Configuration{
	Server:     DefaultServer,
	Zap:        DefaultZap,
	Datasource: DefaultDataSource,
	Email:      DefaultEmail,
	Redis:      DefaultRedis,
}
