package main

import (
	"cloud-platform-system/internal/asynctask"
	"cloud-platform-system/internal/config"
	"cloud-platform-system/internal/handler"
	"cloud-platform-system/internal/svc"
	"context"
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
	server := rest.MustNewServer(c.RestConf, rest.WithCors()) // rest.WithCors(): 用于解决跨域问题
	svcGroup.Add(server)

	// 初始化项目全局属性
	srvCtx := svc.NewServiceContext(c)

	// 初始化异步协程池
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	svcGroup.Add(asynctask.NewMyAsyncTaskPoolServer(cancelCtx, srvCtx, cancelFunc))

	// 注册路由
	handler.RegisterHandlers(server, srvCtx)

	// 启动
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	svcGroup.Start()
}

type PprofServer struct {
	port int
}

// 开启pprof, 便于排除线上问题

func NewPprofServer(port int) *PprofServer {
	return &PprofServer{port: port}
}

func (s *PprofServer) Start() {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
		log.Fatal(err)
	}
	logx.Infof("Start pprof server, listen %d\n", s.port)
}

func (s *PprofServer) Stop() {
	logx.Info("Stop pprof server")
}
