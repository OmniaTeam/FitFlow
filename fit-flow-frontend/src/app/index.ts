import { router, store } from './providers'
import { createApp } from 'vue'
import { Locale } from 'vant'
import ruRU from 'vant/es/locale/lang/ru-RU'

import App from './App.vue'
import 'vant/lib/index.css'

// Устанавливаем русскую локализацию для компонентов Vant
Locale.use('ru-RU', ruRU)

export const app = createApp(App).use(router).use(store)
app.mount('#app')
