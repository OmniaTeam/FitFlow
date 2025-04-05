<template>
	<div class="onboarding-container">
		<template v-if="activeQuestion !== 10">
			<ProgressBar
				:active-question="activeQuestion"
				:total-questions="10"
				@back="decrementActiveQuestion"
			/>
		</template>

		<template v-if="activeQuestion === 1">
			<SetBirthday with-text with-button settings-text="Расскажите нам когда вы родились">
				<template #button>
					<template v-if="userSettings.birthday">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetBirthday>
		</template>
		<template v-if="activeQuestion === 2">
			<SetGender with-text with-button settings-text="Выберите ваш пол">
				<template #button>
					<template v-if="userSettings.gender">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetGender>
		</template>
		<template v-if="activeQuestion === 3">
			<SetWeight with-text with-button settings-text="Выберите ваш текущий вес">
				<template #button>
					<template v-if="userSettings.weight">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetWeight>
		</template>
		<template v-if="activeQuestion === 4">
			<SetHeight with-text with-button settings-text="Выберите свой рост">
				<template #button>
					<template v-if="userSettings.height">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetHeight>
		</template>
		<template v-if="activeQuestion === 5">
			<SetPurpose with-text with-button settings-text="Выберите свою главную цель">
				<template #button>
					<template v-if="userSettings.purpose">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetPurpose>
		</template>
		<template v-if="activeQuestion === 6">
			<SetPlacement with-text with-button settings-text="Где вы планируете заниматься?">
				<template #button>
					<template v-if="userSettings.placement">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetPlacement>
		</template>
		<template v-if="activeQuestion === 7">
			<SetLevel with-text with-button settings-text="Выберите свой уровень подготовки">
				<template #button>
					<template v-if="userSettings.level">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetLevel>
		</template>
		<template v-if="activeQuestion === 8">
			<SetTrainCount with-text with-button settings-text="Укажите, как часто вы планируете заниматься">
				<template #button>
					<template v-if="userSettings.trainCount">
						<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
							background-color: #B22E1A;
							border-color: #B22E1A;
							border-radius: 12px;
							font-weight: bold;
						">
							Дальше
						</Button>
					</template>
				</template>
			</SetTrainCount>
		</template>
		<template v-if="activeQuestion === 9">
			<SetFoodPrompt with-text with-button settings-text="Введите ваши пожелания в составлении питания">
				<template #button>
					<Button @click="incrementActiveQuestion" :type="'primary'" :size="'large'" style="
						background-color: #B22E1A;
						border-color: #B22E1A;
						border-radius: 12px;
						font-weight: bold;
					">
						Дальше
					</Button>
				</template>
			</SetFoodPrompt>
		</template>
		<template v-if="activeQuestion === 10 && userEditState === 'PENDING'">
			<div class="loader">
				<Loading />
			</div>
		</template>
	</div>
</template>

<script lang="ts" setup>
import { useOnBoarding } from '@/widgets/on-boarding/model'
import { Button, Loading } from 'vant'
import { ProgressBar } from '@/shared/ui/progress-bar'
import {
	SetBirthday,
	SetGender,
	SetHeight, SetLevel,
	SetPlacement,
	SetPurpose,
	SetWeight,
	SetTrainCount,
	SetFoodPrompt
} from '@/features/settings'
const {
	userSettings,
	userEditState,
	activeQuestion,

	incrementActiveQuestion,
	decrementActiveQuestion,
} = useOnBoarding(10)
</script>

<style lang="scss" scoped>
.onboarding-container {
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
	gap: 32px;
}

.question {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	gap: 16px;
	height: 100%;

	&__field {
		border-radius: 8px;
		margin-top: 16px;
		background-color: #E6E6E6;
	}
}

.loading {
	position: absolute;
	left: 0;
	right: 0;
	top: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.4);
	display: flex;
	align-items: center;
	justify-content: center;
	flex-direction: column;
}
</style>
