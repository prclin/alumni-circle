package config

type Email struct {
	Host     string
	Port     int
	Username string
	Password string
}

var DefaultEmail = &Email{}
