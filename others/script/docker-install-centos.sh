#!/bin/bash

# 传入一个参数: install/update/delete

# 安装docker
install_docker() {
  # 判断当前系统是否安装了docker
  docker -v
  # shellcheck disable=SC2181
  if [ $? -eq 0 ]; then
    echo "当前机器以安装docker"
    exit
  fi

  # 下载docker
  sudo yum -y update
  sudo yum install -y yum-utils
  sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  sudo yum install -y docker-ce docker-ce-cli containerd.io

  # 配置国内镜像源
  sudo mkdir -p /etc/docker
  # <<: 重定向符, -: 忽略前导制表符(Tab、空格等)
  sudo tee /etc/docker/daemon.json <<-EOF
  {
    "registry-mirrors": ["https://e46bdzxc.mirror.aliyuncs.com"]
  }
EOF

  # 配置docker开机自启动
  sudo systemctl enable docker
  # 启动docker
  sudo systemctl start docker
  # 打印docker版本信息
  echo "docker版本信息如下:"
  docker -v
}

# 卸载docker
delete_docker() {
  sudo systemctl stop docker
  sudo yum remove -y docker-ce \
    docker-ce-cli \
    containerd.io
  sudo rm -rf /etc/systemd/system/docker.service.d
  sudo rm -rf /etc/systemd/system/docker.service
  sudo rm -rf /var/lib/docker
  sudo rm -rf /var/run/docker
  sudo rm -rf /usr/local/docker
  sudo rm -rf /etc/docker
  sudo rm -rf /usr/bin/docker* /usr/bin/containerd* /usr/bin/runc /usr/bin/ctr
}

set -e  # 一旦有异常直接后面的内容不再继续执行
command=$1
if [ "$command" == "install" ]; then
  install_docker
elif [ "$command" == "update" ]; then
  delete_docker
  install_docker
elif [ "$command" == "delete" ]; then
  delete_docker
else
  echo "必加参数是: install/update/delete"
fi
