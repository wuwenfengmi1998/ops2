#!/bin/bash

APP_NAME="OPSYS"
APP_PATH="./build/$APP_NAME"

# 编译应用
echo "编译应用..."
CGO_ENABLED=0 GOOS=linux go build -o $APP_NAME -ldflags="-s -w" main.go

# 创建目录
echo "创建目录..."
sudo mkdir -p $APP_PATH

# 复制文件
echo "复制文件..."
sudo cp $APP_NAME $APP_PATH/
sudo cp -r defConfig $APP_PATH/
sudo cp -r dist $APP_PATH/
