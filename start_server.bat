@echo off
chcp 936 >nul
echo 正在启动问卷系统...

echo 启动后端服务...
start cmd /k "cd backend && go run main.go"

echo 等待后端服务启动...
timeout /t 5 /nobreak >nul

echo 启动前端服务...
start cmd /k "cd frontend/questionnaire-app && npm run dev"

echo 问卷系统已启动
echo 前端访问地址: http://localhost:5173
echo 后端API地址: http://localhost:8080
echo.
echo 按任意键退出此窗口（服务会在后台继续运行）
pause >nul 