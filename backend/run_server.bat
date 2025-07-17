@echo off
chcp 936 >nul
echo 正在启动问卷系统后端服务...

echo 正在启动Go后端服务...
go run main.go

echo 服务已停止
pause 