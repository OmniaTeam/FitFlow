<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field
				:model-value="userSettings?.height ? String(userSettings.height) : ''"
				class="setting-field"
				name="Рост"
				label="Рост (см)"
				@input="(event: any) => handleHeightChange(event.target.value)"
			/>
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field, Popup, Picker } from 'vant'
import { BaseSettings } from '@/shared/ui/base-setting'

defineProps({
	withText: {
		type: Boolean,
		default: false,
	},
	settingsText: {
		type: String,
		default: '',
	},
	withButton: {
		type: Boolean,
		default: false,
	}
})

const {
	userSettings,
	showHeightPicker,

	toggleHeightPicker,
	handleHeightChange,
} = useSettings()

const getHeightColumns = () => {
	return Array.from({ length: 211 }, (_, i) => ({ text: `${i + 40} см`, value: String(i + 40) })).reverse()
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
