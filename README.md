# cloud platform system 后端
## 虚拟Linux
### centos虚拟
```shell
docker pull registry.cn-shanghai.aliyuncs.com/yore/bigdata:7.8.2003_v1

docker run --privileged=true --cap-add SYS_ADMIN -e container=docker -it -p 10022:22 -p 10080:80 -d --name c1 db0f979c4417 /usr/sbin/init
```

## 数据库
### redis
> 拉取镜像并启动
```shell
docker pull redis

docker run -p 6379:6379 --name redis -v /home/redis/data:/data -d redis
```
> 设置密码
```shell
docker exec -it redis redis-cli

config set requirepass 123456
config get requirepass

auth password
```
- 密码: redis

### MongoDB
> 拉取镜像并启动
```shell
docker pull mongo:4.4

docker run -itd --name mongo -v /home/mongo/data:/data/db -p 27017:27017 mongo:4.4
```
> 设置密码与权限
```shell
 docker exec -it mongo mongo admin

db.createUser({ user:'root',pwd:'root',roles:[ { role:'userAdminAnyDatabase', db: 'admin'},'readWriteAnyDatabase']});
```
- 密码: root

#### 设置唯一索引-email
![img.png](img01.png)

## bug
1) 雪花算法生成ID--出现AAAAAAAAAAA=
   原因: 不固定epoch值, 导致每次重启项目这个值都是当前时间值，会导致生成的雪花ID出现重复。
   后面: 自定义分解int64算法保存到byte数组中
2) Token生成和解析出现问题--2024-02-17T00:47:34.123+08:00    error  (/v5/token/refresh - 127.0.0.1:7608) illegal base64 data at input byte 16
   原因: aes加密如果key是16B，那么每次只能对称加密一个16B的块