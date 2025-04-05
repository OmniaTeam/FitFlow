<template>
	<div class="welcome-messages">
		<p :class="{ 'message-visible': messageVisible }">{{ currentMessage }}</p>
	</div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, computed } from 'vue'
import type { StateType } from '@/shared/types'

const props = defineProps<{
	state: StateType
}>()

const messages = [
	'Приятно познакомиться :)',
	'Формируем рекомендации',
	'Почти закончили',
	'Добро пожаловать в Fit FLOW'
]

const errorMessage = 'Произошла системная ошибка'

const currentStep = ref(0)
const messageVisible = ref(false)
const allowNextStep = ref(true)

const currentMessage = computed(() => {
	if (props.state === 'ERROR') {
		return errorMessage
	}

	return messages[currentStep.value]
})

const showMessage = () => {
	messageVisible.value = true
}

const hideMessage = () => {
	messageVisible.value = false
}

const nextStep = () => {
	if (!allowNextStep.value) return

	hideMessage()

	setTimeout(() => {
		if (props.state === 'ERROR') {
			currentStep.value = messages.length - 1
			showMessage()
			return
		}

		// Если статус успешный и мы еще не дошли до последнего сообщения
		if (props.state === 'SUCCESS' && currentStep.value < messages.length - 1) {
			currentStep.value++
			showMessage()

			// Если это последнее сообщение, останавливаем
			if (currentStep.value === messages.length - 1) {
				allowNextStep.value = false
			}
		}
		// Если статус в процессе и мы еще не дошли до предпоследнего сообщения
		else if (props.state === 'PENDING' && currentStep.value < messages.length - 2) {
			currentStep.value++
			showMessage()

			// Запланировать следующий шаг
			if (currentStep.value < messages.length - 2) {
				setTimeout(nextStep, 1500)
			}
		}
		// В остальных случаях просто показываем текущее сообщение
		else {
			showMessage()
		}
	}, 300)
}

// Следим за изменением состояния
watch(() => props.state, (newState) => {
	// Если статус стал успешным и мы находимся на предпоследнем шаге
	if (newState === 'SUCCESS' && currentStep.value === messages.length - 2) {
		// Показываем последнее сообщение
		allowNextStep.value = true
		nextStep()
	}
	// Если произошла ошибка
	else if (newState === 'ERROR') {
		hideMessage()
		setTimeout(() => {
			currentStep.value = 0 // Устанавливаем индекс, но фактически покажется errorMessage
			showMessage()
		}, 300)
	}
})

onMounted(() => {
	// Начинаем показывать сообщения
	setTimeout(() => {
		showMessage()
	}, 300)

	// Запускаем первый переход
	setTimeout(() => {
		nextStep()
	}, 1500)

	// Запускаем последующие переходы
	const stepInterval = setInterval(() => {
		// Если мы уже на последнем или предпоследнем шаге и статус все еще в процессе
		if ((currentStep.value >= messages.length - 2 && props.state === 'PENDING') ||
			(currentStep.value === messages.length - 1) ||
			!allowNextStep.value) {
			clearInterval(stepInterval)
		} else {
			nextStep()
		}
	}, 1500)
})
</script>

<style lang="scss" scoped>
.welcome-messages {
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

	.message-visible {
		opacity: 1;
	}
}
</style>
