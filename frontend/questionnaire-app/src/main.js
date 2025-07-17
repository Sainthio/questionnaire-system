import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'

// 引入Vant
import Vant from 'vant'
import 'vant/lib/index.css'

import App from './App.vue'
import './style.css'

// 引入设备检测工具
import { addDeviceClass } from './utils/device'

// 添加设备类型到HTML元素
addDeviceClass()

const app = createApp(App)

// 使用插件
app.use(createPinia())
app.use(router)
app.use(Vant)

app.mount('#app')
