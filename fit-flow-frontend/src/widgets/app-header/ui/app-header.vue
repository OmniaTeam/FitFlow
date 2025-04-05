<template>
	<header class="app-header">
		<UserIndicator :name="user ? user.firstName ?? '' : ''" :loading="userFetchState === 'PENDING'" />
	</header>
</template>

<script lang="ts" setup>
import { useAppHeader } from '../model'
import { UserIndicator } from '@/shared/ui/user-indicator'

const { userFetchState, user } = useAppHeader()
</script>

<style lang="scss" scoped></style>
