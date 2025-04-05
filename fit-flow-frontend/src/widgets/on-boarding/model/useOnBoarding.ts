import { useStore } from 'vuex'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { IUser } from '@/entities/user'
import type { StateType } from '@/shared/types'

export const useOnBoarding = (lastPoint: number) => {
	const store = useStore()
	const router = useRouter()

	const userSettings = computed<IUser>(() => store.getters['user/getUserSettings'])
	const userEditState = computed<StateType>(() => store.getters['user/getUserEditState'])

	const activeQuestion = ref<number>(1)

	const incrementActiveQuestion = () => {
		activeQuestion.value++
		if (activeQuestion.value === lastPoint) {
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
			}).then(() => {
				if (userEditState.value === 'SUCCESS') {
					router.push('/')
				}
			})
		}
	}
	const decrementActiveQuestion = () => activeQuestion.value--

	return {
		userSettings,
		userEditState,
		activeQuestion,

		incrementActiveQuestion,
		decrementActiveQuestion,
	}
}
