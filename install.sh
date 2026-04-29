#!/bin/bash

# 更新代码
echo "正在拉取最新代码..."
git pull

# 构建前端
echo "正在构建前端..."
cd frontend/ops_vue_js
npm run build
cd ../..

# 安装后端
echo "正在安装后端..."
cd backend/my_work
sudo bash install.sh
cd ../..

echo "安装完成！"
