<template>
	<Button
		@click="loginHandler"
		:type="'primary'"
		:size="'large'"
		style="
			background-color: #B22E1A;
			border-color: #B22E1A;
			border-radius: 12px;
			font-weight: bold;
		"
	>
		Войти через VK ID
	</Button>
</template>

<script lang="ts" setup>
import { useLoginWithVkid } from '../model'
import { Button } from 'vant';

const { userLoginState, loginHandler } = useLoginWithVkid()
</script>
