<template>
	<template v-if="router.name === 'auth' || router.name === 'on-boarding'">
		<slot />
	</template>
	<template v-else>
		<AppHeader />
		<slot />
		<BottomBar />
	</template>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { BottomBar } from '@/widgets/bottom-bar'
import { AppHeader } from '@/widgets/app-header'

const router = useRoute()
</script>
