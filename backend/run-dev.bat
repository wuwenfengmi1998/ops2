@echo off
echo Starting OPS Backend Development Server...
echo.

:: 设置Go环境变量
set CGO_ENABLED=1

:: 检查dist目录
if not exist ".\dist" mkdir dist

:: 运行开发服务器
echo Starting server with CGO enabled for SQLite...
echo Server will be available at: http://127.0.0.1:8080
echo.

go run ./cmd/ops-server/main.go

echo.
echo Server stopped.
pause