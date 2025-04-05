<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field :model-value="userSettings.level ? getFormattedLevel(userSettings.level) : ''" class="setting-field"
				is-link readonly name="Уровень подготовки" label="Уровень подготовки" @click="toggleLevelPicker" />
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup
		v-model:show="showLevelPicker"
		destroy-on-close
		position="bottom"
	>
		<Picker
			:model-value="userSettings.level ? [userSettings.level] : []"
			:columns="getLevelColumns()"
			:default-index="userSettings?.level
				? getLevelColumns().findIndex(item => item.value === userSettings.level)
				: 1"
			@confirm="handleLevelChange"
			@close="showLevelPicker = false"
		/>
	</Popup>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field, Picker, Popup } from 'vant'
import { ELevel } from '@/entities/user/model/enums'
import { getFormattedLevel } from '@/shared/lib/formatters'
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
	showLevelPicker,

	toggleLevelPicker,
	handleLevelChange,
} = useSettings()

const getLevelColumns = () => {
	return [
		{ text: getFormattedLevel(ELevel.pro), value: ELevel.pro },
		{ text: getFormattedLevel(ELevel.medium), value: ELevel.medium },
		{ text: getFormattedLevel(ELevel.new), value: ELevel.new },
	]
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
