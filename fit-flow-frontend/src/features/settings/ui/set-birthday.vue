<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p class="setting__text">{{ settingsText }}</p>
			</template>
			<Field
				:model-value="userSettings.birthday ? getFormattedDate(userSettings.birthday) : ''"
				class="setting-field"
				is-link
				readonly
				name="День рождения"
				label="День рождения"
				@click="toggleBirthdayPicker"
			/>
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup
		v-model:show="showBirthdayPicker"
		destroy-on-close
		position="bottom"
	>
		<DatePicker
			:model-value="userSettings.birthday ? userSettings.birthday.split('T')[0].split('-') : []"
			type="date"
			title="Выберите дату рождения"
			:min-date="new Date(1950, 0, 1)"
			:max-date="new Date()"
			@confirm="handleBirthdayChange"
			@close="showBirthdayPicker = false"
		/>
	</Popup>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { DatePicker, Field, Popup } from 'vant'
import { getFormattedDate } from '@/shared/lib/formatters'
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
	showBirthdayPicker,

	toggleBirthdayPicker,
	handleBirthdayChange,
} = useSettings()
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
