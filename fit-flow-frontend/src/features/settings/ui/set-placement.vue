<template>
	<BaseSettings>
		<template #header>
			<template v-if="withText">
				<p>{{ settingsText }}</p>
			</template>
			<Field :model-value="userSettings.placement ? getFormattedPlacement(userSettings.placement) : ''"
				class="setting-field" is-link readonly name="Место тренировок" label="Место тренировок"
				@click="togglePlacementPicker" />
		</template>
		<template v-if="withButton">
			<slot name="button" />
		</template>
	</BaseSettings>
	<Popup
		v-model:show="showPlacementPicker"
		destroy-on-close
		position="bottom"
	>
		<Picker
			:model-value="userSettings.placement ? [userSettings.placement] : []"
			:columns="getPlacementColumns()"
			:default-index="userSettings?.placement
				? getPlacementColumns().findIndex(item => item.value === userSettings.placement)
				: 0"
			@confirm="handlePlacementChange"
			@close="showPlacementPicker = false"
		/>
	</Popup>
</template>

<script lang="ts" setup>
import { useSettings } from '../model'
import { Field, Picker, Popup } from 'vant'
import { EPlacement } from '@/entities/user/model/enums'
import { getFormattedPlacement } from '@/shared/lib/formatters'
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
	showPlacementPicker,

	togglePlacementPicker,
	handlePlacementChange,
} = useSettings()

const getPlacementColumns = () => {
	return [
		{ text: getFormattedPlacement(EPlacement.home), value: EPlacement.home },
		{ text: getFormattedPlacement(EPlacement.gym), value: EPlacement.gym },
	]
}
</script>

<style lang="scss" scoped>
.setting-field {
	border-radius: 8px;
	background-color: #E6E6E6;
}
</style>
