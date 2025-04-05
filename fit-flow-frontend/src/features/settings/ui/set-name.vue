<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p class="setting__text">{{ settingsText }}</p>
			</template>
			<Field
				:model-value="userSettings?.firstName || ''"
				class="setting-field"
				name="Имя"
				label="Имя"
				@input="(event: any) => handleNameChange(event.target.value)"
			/>
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field } from 'vant'
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
	handleNameChange,
} = useSettings()
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
