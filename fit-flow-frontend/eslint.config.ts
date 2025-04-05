import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import oxlint from 'eslint-plugin-oxlint'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'

export default defineConfigWithVueTs(
	{
		name: 'app/files-to-lint',
		files: ['**/*.{ts,mts,tsx,vue}'],
	},

	{
		name: 'app/files-to-ignore',
		ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
	},

	// Подключаем все vue essential конфиги
	...pluginVue.configs['flat/essential'],

	// Переопределяем нужное правило
	{
		name: 'app/vue-rule-overrides',
		rules: {
			'vue/max-attributes-per-line': ['error', {
				singleline: { max: 1 },
				multiline: { max: 1 },
			}],
		},
	},

	vueTsConfigs.recommended,
	...oxlint.configs['flat/recommended'],
	skipFormatting,
)
