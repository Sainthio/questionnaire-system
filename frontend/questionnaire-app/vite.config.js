import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // 明确指定监听所有网络接口
    port: 5173,
    strictPort: false, // 如果端口被占用则尝试下一个可用端口
    proxy: {
      '/api': {
        target: 'http://0.0.0.0:8080', // 使用0.0.0.0替代localhost
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  }
})
