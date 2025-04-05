<template>
	<div :class="['nav-tab', { 'nav-tab--active': active }]" role="button" tabindex="0" :aria-label="ariaLabel"
		@click="onClick" @keydown.enter="onClick" @keydown.space.prevent="onClick">
		<div class="nav-tab__icon">
			<component :is="iconComponent" :active="active"></component>
		</div>
	</div>
</template>

<script lang="ts" setup>
import type { EmitFn, PropType } from 'vue'
import { useNavTab } from '../model'

const props = defineProps({
	active: {
		type: Boolean,
		required: true,
	},
	icon: {
		type: String as PropType<'home' | 'train' | 'food' | 'chat' | 'account'>,
		required: true,
	},
	ariaLabel: {
		type: String,
		default: 'nav-tab',
	}
})
const emit = defineEmits(['click'])

const { iconComponent, onClick } = useNavTab(props.icon, emit as EmitFn)
</script>

<style lang="scss" scoped>
.nav-tab {
	width: 40px;
	height: 40px;
	display: flex;
	justify-content: center;
	align-items: center;
	border-radius: 10px;
	cursor: pointer;
	outline: none;

	&:focus-visible {
		box-shadow: 0 0 0 2px #007bff;
	}

	&--active {
		background-color: #E6E6E6;
	}

	&__icon {
		width: 24px;
		height: 24px;
		display: flex;
		justify-content: center;
		align-items: center;
	}
}
</style>
