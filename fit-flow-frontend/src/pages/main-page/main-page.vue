<template>
	<div class="main-page">
		main page
	</div>
</template>

<script lang="ts" setup>
</script>

<style lang="scss" scoped>
.main-page {
	display: flex;
	flex-direction: column;
	gap: 20px;

	&__title {
		font-size: 28px;
		font-weight: bold;
		text-align: center;
		color: #E6E6E6;
	}

	&__content {
		display: flex;
		flex-direction: column;
		gap: 20px;
		align-items: center;
		margin-top: 40px;

		p {
			font-size: 18px;
			text-align: center;
			color: #CCCCCC;
		}
	}

	&__auth-link {
		display: inline-block;
		background-color: #B22E1A;
		color: white;
		padding: 12px 24px;
		border-radius: 12px;
		font-weight: bold;
		text-decoration: none;
		transition: background-color 0.3s ease;

		&:hover {
			background-color: #d43b26;
		}
	}
}
</style>
