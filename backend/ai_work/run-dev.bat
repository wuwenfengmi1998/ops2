@echo off
REM OPS Backend 开发服务器启动脚本
REM 设置 CGO 和 GCC 路径
set CGO_ENABLED=1
set PATH=C:\TDM-GCC-64\bin;%PATH%

cd /d %~dp0

echo Starting OPS Backend (with auto-reload)...
echo.

REM 使用 gin 自动重载（比 fresh 更稳定）
go run . -port 8080
