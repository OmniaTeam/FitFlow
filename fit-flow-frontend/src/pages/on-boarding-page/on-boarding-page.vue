<template>
	<div class="onboarding-page">
		<template v-if="showWelcome">
			<div class="welcome-message">
				<h1 :class="{ 'welcome-visible': welcomeVisible }">{{ currentMessage }}</h1>
			</div>
		</template>
		<template v-else>
			<OnBoarding />
		</template>
	</div>
</template>

<script lang="ts" setup>
import { OnBoarding } from '@/widgets/on-boarding'
import { ref, onMounted, computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const userSettings = computed(() => store.getters['user/getUserSettings'])
const userName = computed(() => userSettings.value?.firstName || 'пользователь')

const showWelcome = ref(true)
const welcomeVisible = ref(false)
const currentStep = ref(0)

const currentMessage = computed(() => {
	if (currentStep.value === 0) {
		return `Привет, ${userName.value}👋`
	} else {
		return 'Давай знакомиться'
	}
})

const nextStep = () => {
	welcomeVisible.value = false

	setTimeout(() => {
		if (currentStep.value < 1) {
			currentStep.value++
			welcomeVisible.value = true
		} else {
			showWelcome.value = false
		}
	}, 300)
}

onMounted(() => {
	setTimeout(() => {
		welcomeVisible.value = true
	}, 300)

	setTimeout(() => {
		nextStep()
	}, 1500)

	setTimeout(() => {
		nextStep()
	}, 3000)
})
</script>

<style lang="scss" scoped>
.onboarding-page {
	width: 100%;
	height: 100%;
	position: relative;
}

.welcome-message {
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	width: 100%;
	text-align: center;

	h1 {
		font-size: 28px;
		color: #E6E6E6;
		text-align: center;
		opacity: 0;
		transition: opacity 0.5s ease;
		font-family: 'SF Pro Display', sans-serif;
		font-weight: bold;
	}

	.welcome-visible {
		opacity: 1;
	}
}
</style>
