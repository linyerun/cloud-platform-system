package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Salt    string
	Captcha struct {
		Width  int
		Height int
	}
	Mongo struct {
		Address    string
		Port       int
		Username   string
		Password   string
		AuthSource string
		DbName     string
	}
	Redis struct {
		Address  string
		Port     int
		Password string
	}
	Jwt struct {
		ExpireSec int64
	}
}
