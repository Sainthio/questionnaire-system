@echo off
echo 正在编译重置管理员密码程序...
cd %~dp0
go build -o reset_admin.exe cmd/reset_admin/main.go
if %errorlevel% neq 0 (
    echo 编译失败！
    pause
    exit /b %errorlevel%
)

echo 正在重置管理员密码...
reset_admin.exe
if %errorlevel% neq 0 (
    echo 重置失败！
    pause
    exit /b %errorlevel%
)

echo.
echo 管理员密码重置成功！
echo 用户名: admin
echo 密码: admin123
echo.
pause 