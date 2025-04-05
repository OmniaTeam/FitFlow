<template>
	<div class="main-settings">
		<p class="main-settings__path-title">Основные настройки</p>
		<div class="main-settings__fields">
			<SetName />
			<SetBirthday />
			<SetGender />
			<SetHeight />
			<SetWeight />
		</div>
	</div>
</template>

<script lang="ts" setup>
import { SetName, SetBirthday, SetGender, SetHeight, SetWeight } from '@/features/settings'
</script>

<style lang="scss" scoped>
.main-settings {
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
