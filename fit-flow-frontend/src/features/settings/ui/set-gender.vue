<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field :model-value="userSettings.gender ? getFormattedGender(userSettings.gender) : ''"
				class="setting-field" is-link readonly name="Пол" label="Пол" @click="toggleGenderPicker" />
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup
		v-model:show="showGenderPicker"
		destroy-on-close
		position="bottom"
	>
		<Picker
			:model-value="userSettings.gender ? [userSettings.gender] : []"
			:columns="getGenderColumns()"
			:default-index="userSettings?.gender
				? getGenderColumns().findIndex(item => item.value === userSettings.gender)
				: 0"
			@confirm="handleGenderChange"
			@close="showGenderPicker = false"
		/>
	</Popup>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field, Popup, Picker } from 'vant'
import { EGender } from '@/entities/user/model/enums'
import { getFormattedGender } from '@/shared/lib/formatters'
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
	showGenderPicker,

	toggleGenderPicker,
	handleGenderChange,
} = useSettings()

const getGenderColumns = () => {
	return [
		{ text: getFormattedGender(EGender.male), value: EGender.male },
		{ text: getFormattedGender(EGender.female), value: EGender.female }
	]
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
