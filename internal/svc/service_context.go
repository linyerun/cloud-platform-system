package svc

import (
	"cloud-platform-system/internal/config"
	"cloud-platform-system/internal/middleware"
	"cloud-platform-system/internal/models"
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext struct {
	// 项目配置
	Config config.Config

	// 端口管理员
	PortManager *portObj

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
		Config: c,

		PortManager: newPortObj(redisClient, mongoClient, c),

		RedisClient:  redisClient,
		MongoClient:  mongoClient,
		DockerClient: client,

		CAPTCHA:       "Captcha:%v",
		ExceptionList: "ExceptionJsonList",
		ImagePrefix:   "image_%v",

		JwtAuth: middleware.NewJwtAuthMiddleware().Handle,
		Visitor: middleware.NewVisitorMiddleware().Handle,
		User:    middleware.NewUserMiddleware().Handle,
		Admin:   middleware.NewAdminMiddleware().Handle,
		Super:   middleware.NewSuperMiddleware().Handle,
	}
}

type portObj struct {
	SinglePortKey string
	TenPortKey    string
	RedisClient   *redis.Client
	MongoClient   *mongo.Client
	Config        config.Config
}

func newPortObj(redisCli *redis.Client, mongoCli *mongo.Client, config config.Config) *portObj {
	singlePortKey := "single_port_idx"
	tenPortKey := "ten_port_idx"
	// 如果redis中不存在single_port_idx, 初始化为10000
	err := redisCli.SetNX(context.Background(), singlePortKey, 10000, redis.KeepTTL).Err()
	if err != nil {
		panic(err)
	}
	// 如果redis中不存在ten_port_idx, 初始化为65535
	err = redisCli.SetNX(context.Background(), tenPortKey, 65535, redis.KeepTTL).Err()
	if err != nil {
		panic(err)
	}
	return &portObj{RedisClient: redisCli, MongoClient: mongoCli, SinglePortKey: singlePortKey, TenPortKey: tenPortKey, Config: config}
}

func (o *portObj) GetSinglePort() (port uint, err error) {
	// 从Mongo中查看是否存在能复用的
	one := o.MongoClient.Database(o.Config.Mongo.DbName).Collection(models.PortRecycleDocument).FindOneAndDelete(context.Background(), bson.D{{"type", models.SinglePortType}}, options.FindOneAndDelete().SetProjection(bson.D{{"port_start", 1}}))
	if err = one.Err(); err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	} else if err == nil {
		portRecycle := new(models.PortRecycle)
		err = one.Decode(portRecycle)
		if err != nil {
			return 0, err
		} else {
			return portRecycle.PortStart, nil
		}
	}

	// 判断是否还能再获取一个端口
	singlePortIdx, tenPortIdx, err := o.getIdx()
	if err != nil {
		return 0, err
	}
	if tenPortIdx < singlePortIdx {
		return 0, errors.New("not one ports")
	}

	// 从redis中获取一个端口并更新相对应的idx
	err = o.RedisClient.Set(context.Background(), o.SinglePortKey, singlePortIdx+1, redis.KeepTTL).Err()
	if err != nil {
		return 0, err
	}
	return uint(singlePortIdx), nil
}

func (o *portObj) GetTenPort() (from, to uint, err error) {
	// 从Mongo中查看是否存在能复用的
	one := o.MongoClient.Database(o.Config.Mongo.DbName).Collection(models.PortRecycleDocument).FindOneAndDelete(context.Background(), bson.D{{"type", models.TenPortsType}}, options.FindOneAndDelete().SetProjection(bson.D{{"port_start", 1}}))
	if err = one.Err(); err != nil && err != mongo.ErrNoDocuments {
		return 0, 0, err
	} else if err == nil {
		portRecycle := new(models.PortRecycle)
		err = one.Decode(portRecycle)
		if err != nil {
			return 0, 0, err
		} else {
			return portRecycle.PortStart, portRecycle.PortStart + 9, nil
		}
	}

	// 判断是否还能再获取十个端口
	singlePortIdx, tenPortIdx, err := o.getIdx()
	if tenPortIdx-singlePortIdx+1 < 10 {
		return 0, 0, errors.New("not ten ports")
	}

	// 从redis中获取十个端口并更新相对应的idx
	err = o.RedisClient.Set(context.Background(), o.TenPortKey, tenPortIdx-10, redis.KeepTTL).Err()
	if err != nil {
		return 0, 0, err
	}
	return uint(tenPortIdx) - 9, uint(tenPortIdx), nil
}

func (o *portObj) RepaySinglePort(port uint) error {
	// 把结果保存到Mongo中
	portRecycle := models.PortRecycle{Id: fmt.Sprintf("%d-%d", port, port), Type: models.SinglePortType, PortStart: port}
	_, err := o.MongoClient.Database(o.Config.Mongo.DbName).Collection(models.PortRecycleDocument).InsertOne(context.Background(), portRecycle)
	return err
}

func (o *portObj) RepayTenPort(from, to uint) error {
	if to-from+1 != 10 {
		return errors.New("from到to之间不是10个端口")
	}
	// 把结果保存到Mongo中
	portRecycle := models.PortRecycle{Id: fmt.Sprintf("%d-%d", from, to), Type: models.TenPortsType, PortStart: from}
	_, err := o.MongoClient.Database(o.Config.Mongo.DbName).Collection(models.PortRecycleDocument).InsertOne(context.Background(), portRecycle)
	return err
}

func (o *portObj) getIdx() (singlePortIdx, tenPortIdx int, err error) {
	singlePortIdx, err = o.RedisClient.Get(context.Background(), o.SinglePortKey).Int()
	if err != nil {
		return 0, 0, err
	}
	tenPortIdx, err = o.RedisClient.Get(context.Background(), o.TenPortKey).Int()
	if err != nil {
		return 0, 0, err
	}
	return
}
