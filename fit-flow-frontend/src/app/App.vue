<template>
	<main class="main-application">
		<AuthProvider>
			<AppLayout>
				<RouterView />
			</AppLayout>
		</AuthProvider>
	</main>
</template>

<script setup lang="ts">
import { AuthProvider } from '@/app/providers/auth'
import { AppLayout } from '@/app/layouts/app-layout'

import './styles/index.css'

</script>

<style scoped lang="scss">
.main-application {
	overflow-x: hidden;
	max-width: 490px;
	padding: 20px 20px 140px 20px;
	margin: 0 auto;
	height: calc(100vh - 160px);
}
</style>
