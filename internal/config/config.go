package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Salt    string
	Captcha struct {
		Width      int
		Height     int
		TimeoutSec uint
	}
	Admin struct {
		Email string
	}
	PortManager struct {
		From uint
		To   uint
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
	Pprof struct {
		Port int
	}
	AsyncTask struct {
		PullTaskWaitMillSec int64
	}
	Container struct {
		Host         string
		InitUsername string
		InitPassword string
	}
}
