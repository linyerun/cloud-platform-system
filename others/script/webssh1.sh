#!/bin/bash

install_webssh() {
  # 拉取镜像
  docker pull registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh
  echo "===>成功拉取到webssh镜像"

  # 运行镜像
  docker run --name webssh -d -p 8888:8888 registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh
  echo "===>成功启动了webssh镜像"

  # 打印信息
  echo "===>访问的页面地址: http://服务器IP:8888"
}

delete_webssh() {
  image_id=$(docker images -a | grep "registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh" | awk '{print $3}')
  container_id=$(docker ps -a | grep -E "$image_id|registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh" | awk '{print $1}')

  # 关闭并删除容器
  docker stop "$container_id"
  docker rm "$container_id"
  echo "===>关闭${container_id}成功"

  # 删除镜像
  docker rmi "$image_id"
  echo "===>删除${image_id}成功"
}

update_webssh() {
  delete_webssh
  install_webssh
}

set -e  # 一旦有异常直接后面的内容不再继续执行
command=$1
if [ "$command" == "install" ]; then
  install_webssh
elif [ "$command" == "delete" ]; then
  delete_webssh
elif [ "$command" == "update" ]; then
  update_webssh
else
  echo "===>必须带的参数为: install/update/delete"
fi
