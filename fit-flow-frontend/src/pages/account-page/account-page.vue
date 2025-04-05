<template>
	<div class="account-page">
		<MainSettings />
		<CommonSettings />
		<div class="account-page__buttons">
			<Button
				@click="saveHandler"
				:type="'primary'"
				:size="'large'"
				:loading="userEditState === 'PENDING'"
				style="
					background-color: #B22E1A;
					border-color: #B22E1A;
					border-radius: 12px;
					font-weight: bold;
				"
			>
				Сохранить настройки
			</Button>
			<Button
				@click="logoutHandler"
				:type="'primary'"
				:size="'large'"
				:loading="userLogoutState === 'PENDING'"
				style="
					background-color: #B22E1A;
					border-color: #B22E1A;
					border-radius: 12px;
					font-weight: bold;
				"
			>
				Выйти из аккаунта
			</Button>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { MainSettings } from '@/widgets/main-settings'
import { CommonSettings } from '@/widgets/common-settings'
import { Button } from 'vant'
import { useStore } from 'vuex'
import { computed } from 'vue'

const store = useStore()

const userSettings = computed(() => store.getters['user/getUserSettings'])
const userEditState	= computed(() => store.getters['user/getUserEditState'])
const userLogoutState= computed(() => store.getters['user/getuserLogoutState'])

const saveHandler = () => {
	store.dispatch('user/editUser', {
		birthday: userSettings.value.birthday,
		firstName: userSettings.value.firstName,
		foodPrompt: userSettings.value.foodPrompt,
		gender: userSettings.value.gender,
		height: userSettings.value.height,
		lastName: userSettings.value.lastName,
		level: userSettings.value.level,
		placement: userSettings.value.placement,
		purpose: userSettings.value.purpose,
		trainingCount: userSettings.value.trainCount,
		weight: userSettings.value.weight,
	})
}

const logoutHandler = () => {
	store.dispatch('user/logout')
}
</script>

<style lang="scss" scoped>
.account-page {
	padding-top: 32px;
	display: flex;
	flex-direction: column;
	gap: 16px;
	&__buttons {
		margin-top: 32px;
		display: flex;
		flex-direction: column;
		gap: 16px;
	}
}
</style>
