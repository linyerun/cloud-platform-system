#!/bin/bash


install_webssh2() {
  # 拉取镜像
  docker pull registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh2
  echo "===>成功拉取到webssh2镜像"

  # 启动镜像
  docker run --name webssh2 -d -p 2222:2222 registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh2
  echo "===>成功运行webssh2镜像"

  # 打印使用镜像的信息
  echo "===>访问的页面地址: http://服务器IP:2222/ssh/host/想连的主机ip"
  echo "===>之后输入相连的主机的用户名和密码"
}

delete_webssh2() {
  image_id=$(docker images -a | grep "registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh2" | awk '{print $3}')
  container_id=$(docker ps -a | grep -E "$image_id|registry.cn-hangzhou.aliyuncs.com/lyr_public/webssh2" | awk '{print $1}')

  # 关闭并删除容器
  docker stop "$container_id"
  docker rm "$container_id"
  echo "===>关闭${container_id}成功"

  # 删除镜像
  docker rmi "$image_id"
  echo "===>删除${image_id}成功"
}

update_webssh2() {
  delete_webssh2
  install_webssh2
}

set -e  # 一旦有异常直接后面的内容不再继续执行
command=$1
if [ "$command" == "install" ]; then
  install_webssh2
elif [ "$command" == "delete" ]; then
  delete_webssh2
elif [ "$command" == "update" ]; then
  update_webssh2
else
  echo "===>必须带的参数为: install/update/delete"
fi
