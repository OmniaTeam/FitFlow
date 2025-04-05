<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field
				:model-value="userSettings.foodPrompt || ''"
				class="setting-field"
				type="textarea"
				maxlength="500"
				autosize
				@input="(event: any) => handleFoodPromptChange(event.target.value)"
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
	handleFoodPromptChange,
} = useSettings()
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
