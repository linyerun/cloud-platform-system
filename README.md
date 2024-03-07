# cloud platform system 后端
## 虚拟Linux
### centos虚拟
#### 指令
```shell
docker pull registry.cn-shanghai.aliyuncs.com/yore/bigdata:7.8.2003_v1

docker run --privileged=true --cap-add SYS_ADMIN -e container=docker -it -p 10022:22 -p 10080:80 -d --name c1 db0f979c4417 /usr/sbin/init
```

> 新版本: registry.cn-hangzhou.aliyuncs.com/lyr_public/centos7:2.0
- 运行指令: `docker run --privileged=true -d -p 10022:22 -p 10080:80 --name c1 registry.cn-hangzhou.aliyuncs.com/lyr_public/centos7:2.0`
- 账户：root，密码：root

#### 端口说明
- Linux的端口范围: 0-65535
- 留给服务器自己使用的端口范围: 0-9999
- 给用户部署项目/开启数据库使用的端口：10000开始
- 给Linux虚拟服务器分配端口：从65535往左边靠近来使用，每次分配10个端口

### debian
- 下载东西：https://blog.csdn.net/u012285295/article/details/134308940
- 配置ssh：https://blog.csdn.net/future_ai/article/details/81701744

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
3) debian执行`apt update`失败，解决方法是更改时区并修改时间。https://blog.csdn.net/weixin_45663954/article/details/123394680

### 修理debian
1) debian执行`apt update`失败，解决方法是更改时区并修改时间。https://blog.csdn.net/weixin_45663954/article/details/123394680
2) 修改镜像源
   1. apt install apt-transport-https ca-certificates
   2. 配合docker cp和本身的mv指令完成镜像源的修改
   3. 执行`apt-get update && apt update`
   4. 安装 `apt-get install -y vim`
```text
# 默认注释了源码镜像以提高 apt update 速度
# 清华大学的软件源
deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free 
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free 
deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free 
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free 
deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free 
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free 
deb https://mirrors.tuna.tsinghua.edu.cn/debian-security bullseye-security main contrib non-free 
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian-security bullseye-security main contrib non-free
```
## ubuntu

### 无法连接上ssh

- 使用了`systemctl start sshd && ufw disable`还是不行
  1.  `vim /etc/ssh/ssh_config`
  2. `PermitRootLogin yes`，`PasswordAuthentication yes`

## git指令

```text
安装git后
$ git config --global user.name "Your Name"
$ git config --global user.email "email@example.com"

$ ssh-keygen -t rsa -C "youremail@example.com" 创建ssh key，用于和github通信
(秘钥存储于C:\Users\27634\.ssh，把公钥id_rsa.pub存储于github)

创建版本库
$ pwd 命令用于显示当前目录(没啥用)
$ git init 把这个目录变成Git可以管理的仓库(后续新建提交和ssh克隆需要)	

操作版本库
$ git add 文件名 添加文件(新增或者更改都需要先add)
$ git commit -m "说明" 提交到本地版本库

$ git status 查看仓库状态
$ git diff 文件名 查看修改的地方

版本回退(从一个commit恢复)
$ git log 查看版本历史
$ git reset --hard HEAD^ 回退到上个版本
$ git reset --hard 1094a 回退到特定版本号(commit以后回退)
$ git reflog 记录每一次命令

$ git checkout -- file 直接丢弃工作区的修改(add以前回退)
$ git reset HEAD <file> 添加到了暂存区时，想丢弃修改(add以后回退)

删除文件
$ git rm file(已经add/commit,在目录中删除)

$ git checkout -- file 删错了回退

远程仓库
$ git remote add origin git@server-name:path/repo-name.git 关联远程库
$ git push -u origin master 第一次的push
$ git push origin master 常用的push，本地分支会在服务器上新建分支
$ git pull 需要有关联的分支，第一次下拉最好新建一个空文件夹
$ git branch --set-upstream-to=origin/远程分支 本地分支 关联分支

$ git clone git@server-name:path/repo-name.git 克隆(不需要另建文件夹)

分支
$ git branch -a 查看所有分支
$ git branch -vv 查看分支关联
$ git branch dev 创建分支
$ git checkout dev 切换分支
$ git merge dev 合并某分支到当前分支
$ git merge --no-ff -m "msg" dev 普通模式合并，合并后的历史有分支
$ git branch -d dev 删除分支
$ git checkout -b dev 创建并切换分支


合并分支,无法merge
$ git stash save 名字 暂存工作状态
$ git pull origin dev 拉下来 
$ git stash list 查看已经暂存的状态
$ git stash pop stash@{0} 将暂存状态merge到当前分支
还有冲突时,手动修改文件,然后add/commit
$ git log --graph 分支合并图

bug分支issue
$ git stash 暂存工作状态
$ git stash list 查看暂存工作状态
$ git stash pop 恢复暂存状态并删除状态

开发分支feature
$ git branch -D <name> 强制删除未合并的分支

rebase
$ git rebase 本地未push的分叉提交历史整理成直线

标签
$ git tag 标签名 打在最新提交的commit上
$ git tag 查询所有标签
$ git tag 标签名 f52c633 给特定的commit打标签
$ git tag -a 标签名 -m "msg" commit的id 给标签设置说明
$ git show 标签名 查询标签内容
$ git tag -d 标签名 删除标签
$ git push origin 标签名 推送某个标签到远程
$ git push origin --tags 推送所有标签
$ git push origin :refs/tags/<tagname> 可以删除一个远程标签。
```

## 三大镜像总结

| 镜像名称 |                           镜像地址                           | 镜像运行指令 | 已安装软件       | 查看版本       |
| -------- | :----------------------------------------------------------: | ------------ | ---------------- | -------------- |
| centos7  |                                                              |              | systemd, ssh,vim |                |
| debian11 |                                                              |              | systemd, ssh,vim |                |
| ubuntu20 | docker run -itd --privileged=true --name c1 -p 20022:22 registry.cn-hangzhou.aliyuncs.com/lyr_public/ubuntu:2.0 /bin/bash |              | systemd, ssh,vim | lsb_release -a |

> 解决debian11和ubuntu20开启自启动ssh服务问题，然后以后这两个配成`CMD ["systemd"]`

# 设置Linux时间
- date -s "YYYY-YY-DD HH:mm:ss"
  - 手动修改时间
- hwclock -w