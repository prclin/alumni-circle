package config

/*
Configuration 对应整个application.yaml配置文件
*/
type Configuration struct {
	Server     *Server
	Zap        *Zap
	Datasource *Datasource
	Redis      *Redis
	Jwt        *Jwt
}

var DefaultConfiguration = &Configuration{
	Server:     DefaultServer,
	Zap:        DefaultZap,
	Datasource: DefaultDataSource,
	Redis:      DefaultRedis,
	Jwt:        DefaultJwt,
}
