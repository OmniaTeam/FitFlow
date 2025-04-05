<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field
				:model-value="userSettings?.weight ? String(userSettings.weight) : ''"
				class="setting-field"
				name="Вес"
				label="Вес (кг)"
				type="number"
				:max="200"
				@input="(event: any) => handleWeightChange(event.target.value)"
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
	showWeightPicker,

	toggleWeightPicker,
	handleWeightChange,
} = useSettings()

const getWeightColumns = () => {
	return Array.from({ length: 191 }, (_, i) => ({ text: `${i + 10} кг`, value: String(i + 10) })).reverse()
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
