<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field :model-value="userSettings.trainCount ? String(userSettings.trainCount) : ''" class="setting-field"
				is-link readonly name="Количество" label="Количество" @click="toggleTrainCountPicker" />
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup v-model:show="showTrainCountPicker" destroy-on-close position="bottom">
		<Picker :columns="getTrainCountColumns()" :default-index="userSettings?.trainCount
			? getTrainCountColumns().findIndex(item => item.value === String(userSettings.trainCount))
			: Math.floor(7 / 2)" @confirm="handleTrainCountChange" @close="showTrainCountPicker = false" />
	</Popup>
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
	showTrainCountPicker,

	toggleTrainCountPicker,
	handleTrainCountChange,
} = useSettings()

const getTrainCountColumns = () => {
	return Array.from({ length: 7 }, (_, i) => ({ text: `${i + 1} раз`, value: String(i + 1) })).reverse()
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
