import { router, store } from './providers'
import { createApp } from 'vue'

import App from './App.vue'

export const app = createApp(App).use(router).use(store)
app.mount('#app')
