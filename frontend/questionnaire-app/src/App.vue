<script setup>
// App.vue
import { onMounted } from 'vue';
import { isDesktop } from './utils/device';

onMounted(() => {
  // 根据窗口宽度设置适当的类名
  const setDeviceClass = () => {
    if (window.innerWidth > 768) {
      document.documentElement.classList.remove('mobile');
      document.documentElement.classList.add('desktop');
    } else {
      document.documentElement.classList.remove('desktop');
      document.documentElement.classList.add('mobile');
    }
  };
  
  // 初始设置
  setDeviceClass();
  
  // 监听窗口大小变化
  window.addEventListener('resize', setDeviceClass);
});
</script>

<template>
  <div class="app-container">
    <div class="app-content">
      <router-view />
    </div>
  </div>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

.app-container {
  max-width: 100vw;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  background-color: #f7f8fa;
}

.app-content {
  width: 100%;
  max-width: 600px; /* 在桌面端限制最大宽度 */
  min-height: 100vh;
  background-color: #fff;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
  position: relative;
}

/* 桌面端样式 */
.desktop .app-container {
  padding: 20px 0;
  background-color: #f0f2f5;
}

.desktop .app-content {
  max-width: 600px;
  min-height: calc(100vh - 40px);
  border-radius: 8px;
  overflow: hidden;
}

/* 移动端样式 */
.mobile .app-content {
  max-width: 100%;
  box-shadow: none;
  border-radius: 0;
}
</style>
