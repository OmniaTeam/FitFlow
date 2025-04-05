import { fileURLToPath, URL } from 'node:url'
import { VantResolver } from '@vant/auto-import-resolver'
import { defineConfig, loadEnv } from 'vite'

import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
	// Загружаем переменные окружения в зависимости от режима
	const env = loadEnv(mode, process.cwd(), '')

	console.log(`Режим сборки: ${mode}, VITE_MODE: ${env.VITE_MODE}`)

	return {
		plugins: [
			vue(),
			vueJsx(),
			vueDevTools(),
			AutoImport({
				resolvers: [VantResolver()],
			}),
			Components({
				resolvers: [VantResolver()],
			}),
		],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url)),
			},
		},
		// Добавляем определение переменной окружения для клиентского кода
		define: {
			'process.env.NODE_ENV': JSON.stringify(mode),
		},
	}
})
