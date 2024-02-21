package main

import (
	"cloud-platform-system/internal/config"
	"cloud-platform-system/internal/handler"
	"cloud-platform-system/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"log"
	"net/http"
	_ "net/http/pprof" // 目的: 让pprof包执行init函数将请求处理函数注册到default请求处理器中

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/cloud_platform.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcGroup := service.NewServiceGroup()
	defer svcGroup.Stop()

	// 添加pprof
	svcGroup.Add(NewPprofServer(c.Pprof.Port))

	// 添加项目本身的请求处理器
	server := rest.MustNewServer(c.RestConf)
	svcGroup.Add(server)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 启动
	svcGroup.Start()
}

type pprofServer struct {
	port int
}

func NewPprofServer(port int) *pprofServer {
	return &pprofServer{port: port}
}

func (s *pprofServer) Start() {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
		log.Fatal(err)
	}
	logx.Infof("Start pprof server, listen %d\n", s.port)
}

func (s *pprofServer) Stop() {
	logx.Info("Stop pprof server")
}
