#!/bin/bash

APP_NAME="OPSYS"
APP_PATH="/opt/$APP_NAME"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"
LOG_PATH="/var/log/$APP_NAME"

echo "正在安装 $APP_NAME..."

# 编译应用
echo "编译应用..."
CGO_ENABLED=0 GOOS=linux go build -o $APP_NAME -ldflags="-s -w" main.go

# 创建目录
echo "创建目录..."
sudo mkdir -p $APP_PATH
sudo mkdir -p $LOG_PATH

# 复制文件
echo "复制文件..."
sudo cp $APP_NAME $APP_PATH/
sudo cp -r defConfig $APP_PATH/
sudo cp -r dist $APP_PATH/

# 创建服务文件
echo "创建服务文件..."
sudo tee $SERVICE_FILE > /dev/null <<EOF
[Unit]
Description=My Gin Application
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=$APP_PATH
ExecStart=$APP_PATH/$APP_NAME
Restart=always
RestartSec=5
StandardOutput=append:$LOG_PATH/access.log
StandardError=append:$LOG_PATH/error.log
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

# 设置权限
echo "设置权限..."
sudo chown -R www-data:www-data $APP_PATH
sudo chown -R www-data:www-data $LOG_PATH
sudo chmod 750 $APP_PATH/$APP_NAME

# 重载并启动
echo "启动服务..."
sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
sudo systemctl start $APP_NAME

echo "安装完成！"
echo "使用以下命令管理服务："
echo "  sudo systemctl status $APP_NAME"
echo "  sudo journalctl -u $APP_NAME -f"