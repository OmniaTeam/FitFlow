<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field :model-value="userSettings.purpose ? getFormattedPurpose(userSettings.purpose) : ''"
				class="setting-field" is-link readonly name="Цель" label="Цель" @click="togglePurposePicker" />
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup
		v-model:show="showPurposePicker"
		destroy-on-close
		position="bottom"
	>
		<Picker
			:model-value="userSettings.purpose ? [userSettings.purpose] : []"
			:columns="getPurposeColumns()"
			:default-index="userSettings?.purpose
				? getPurposeColumns().findIndex(item => item.value === userSettings.purpose)
				: 1"
			@confirm="handlePurposeChange"
			@close="showPurposePicker = false"
		/>
	</Popup>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field, Picker, Popup } from 'vant'
import { EPurpose } from '@/entities/user/model/enums'
import { getFormattedPurpose } from '@/shared/lib/formatters'
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
	showPurposePicker,

	togglePurposePicker,
	handlePurposeChange,
} = useSettings()

const getPurposeColumns = () => {
	return [
		{ text: getFormattedPurpose(EPurpose.gain), value: EPurpose.gain },
		{ text: getFormattedPurpose(EPurpose.slim), value: EPurpose.slim },
		{ text: getFormattedPurpose(EPurpose.stable), value: EPurpose.stable }
	]
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
