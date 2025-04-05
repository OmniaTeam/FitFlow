<template>
	<div class="common-settings">
		<p class="common-settings__path-title">Общие настройки</p>
		<div class="common-settings__fields">
			<SetPurpose />
			<SetPlacement />
			<SetLevel />
			<SetTrainCount />
			<SetFoodPrompt />
		</div>
	</div>
</template>

<script lang="ts" setup>
import {
	SetPurpose,
	SetPlacement,
	SetLevel,
	SetTrainCount,
	SetFoodPrompt
} from '@/features/settings'
</script>

<style lang="scss" scoped>
.common-settings {
	&__path-title {
		font-size: 20px;
		font-weight: bold;
		color: #E6E6E6;
	}

	&__fields {
		padding-top: 20px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
}
</style>
