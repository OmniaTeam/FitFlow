<template>
	<slot />
</template>

<script lang="ts" setup>
import { useAuthProvider } from '../model'
import { onMounted } from 'vue'

const { checkUserAuth } = useAuthProvider()

onMounted(() => {
	checkUserAuth()
})
</script>
