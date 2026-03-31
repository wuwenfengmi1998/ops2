@echo off
echo Starting OPS backend server (refactored version)...

REM 检查前端dist目录是否存在
if not exist "./dist" (
    echo WARNING: Frontend build not found at ./dist
    echo Please build frontend first or copy build files to ./dist
    echo.
)

REM 运行新的重构版本
echo Running new refactored backend...
echo.
go run ./cmd/ops-server/main.go

pause