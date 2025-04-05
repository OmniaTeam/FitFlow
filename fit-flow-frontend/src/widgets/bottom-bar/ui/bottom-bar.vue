<template>
	<div class="bottom-bar">
		<NavTab :active="activeTab === 'home'" icon="home" @click="changeTab('/')" />
		<NavTab :active="activeTab === 'train'" icon="train" @click="changeTab('train')" />
		<NavTab :active="activeTab === 'food'" icon="food" @click="changeTab('food')" />
		<NavTab :active="activeTab === 'chat'" icon="chat" @click="changeTab('chat')" />
		<NavTab :active="activeTab === 'account'" icon="account" @click="changeTab('account')" />
	</div>
</template>

<script lang="ts" setup>
import { useBottomBar } from '../model'
import { NavTab } from '@/shared/ui/nav-tab'

const { activeTab, changeTab } = useBottomBar()
</script>

<style lang="scss" scoped>
.bottom-bar {
	position: fixed;
	bottom: 32px;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 12px 8px;
	width: calc(100% - 56px);
	max-width: calc(490px - 16px);
	background-color: #010101;
	border-radius: 12px;
}
</style>
