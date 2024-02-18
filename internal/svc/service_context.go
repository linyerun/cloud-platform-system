package svc

import (
	"cloud-platform-system/internal/config"
	"cloud-platform-system/internal/middleware"
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext struct {
	// 项目配置
	Config config.Config

	// 数据库客户端
	RedisClient  *redis.Client
	MongoClient  *mongo.Client
	DockerClient *docker.Client

	// 当常量用的值
	CAPTCHA       string
	ExceptionList string
	ImagePrefix   string

	// 中间件
	JwtAuth rest.Middleware
	Visitor rest.Middleware
	User    rest.Middleware
	Admin   rest.Middleware
	Super   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MongoDB客户端
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", c.Mongo.Address, c.Mongo.Port)).SetAuth(options.Credential{
		Username:   c.Mongo.Username,
		Password:   c.Mongo.Password,
		AuthSource: c.Mongo.AuthSource,
	}))
	if err != nil {
		panic(err)
	}
	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		panic(err)
	}

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Address, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       0,
	})
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	// 初始化Docker
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	// 初始化ServiceContext
	return &ServiceContext{
		Config:        c,
		RedisClient:   redisClient,
		MongoClient:   mongoClient,
		DockerClient:  client,
		CAPTCHA:       "Captcha:%v",
		ExceptionList: "ExceptionJsonList",
		ImagePrefix:   "image_%v",
		JwtAuth:       middleware.NewJwtAuthMiddleware().Handle,
		Visitor:       middleware.NewVisitorMiddleware().Handle,
		User:          middleware.NewUserMiddleware().Handle,
		Admin:         middleware.NewAdminMiddleware().Handle,
		Super:         middleware.NewSuperMiddleware().Handle,
	}
}
