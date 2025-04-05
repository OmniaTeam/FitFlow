<template>
	<div class="user-indicator">
		<template v-if="loading">
			<Skeleton avatar avatar-shape="round" avatar-size="48" style="border-radius: 32px;" />
		</template>
		<template v-else>
			<Image width="48" height="48" round :src="imageUrl" />
		</template>
		<div class="user-info">
			<p class="user-info__hello">Добро пожаловать</p>
			<template v-if="loading">
				<Skeleton title />
			</template>
			<template v-else>
				<p class="user-info__name">
					{{ name }}👋
				</p>
			</template>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { Image } from 'vant'
import { Skeleton } from 'vant/es'

defineProps({
	name: {
		type: String,
		default: "",
	},
	imageUrl: {
		type: String,
		default: "",
	},
	loading: {
		type: Boolean,
		default: false,
	}
})
</script>

<style lang="scss" scoped>
.user-indicator {
	display: flex;
	align-items: center;
	gap: 8px;
}

.user-info {
	display: flex;
	flex-direction: column;
	gap: 4px;

	&__hello {
		color: #A2A2A2;
		font-size: 16px;
	}

	&__name {
		color: #E6E6E6;
		font-size: 20px;
		font-weight: bold;
	}
}
</style>
