/**
 * 设备检测工具
 */

// 检测是否为移动设备
export const isMobile = () => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) || 
         window.innerWidth <= 768;
};

// 检测是否为桌面设备
export const isDesktop = () => {
  return !isMobile();
};

// 获取当前设备类型
export const getDeviceType = () => {
  return isMobile() ? 'mobile' : 'desktop';
};

// 添加设备类型到HTML元素
export const addDeviceClass = () => {
  const deviceType = getDeviceType();
  document.documentElement.classList.add(deviceType);
  
  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    const newDeviceType = getDeviceType();
    if (newDeviceType !== deviceType) {
      document.documentElement.classList.remove(deviceType);
      document.documentElement.classList.add(newDeviceType);
    }
  });
};

export default {
  isMobile,
  isDesktop,
  getDeviceType,
  addDeviceClass
}; 