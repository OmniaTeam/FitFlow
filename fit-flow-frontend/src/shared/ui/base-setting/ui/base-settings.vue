<template>
	<div class="settings">
		<div class="settings__header">
			<slot name="header" />
		</div>
		<slot />
	</div>
</template>

<script lang="ts" setup>
</script>

<style lang="scss" scoped>
.settings {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	gap: 16px;
	height: 100%;
	&__header {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}
}
</style>
