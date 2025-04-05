import { useStore } from 'vuex'
import { computed } from 'vue'
import type { IUser } from '@/entities/user'
import type { StateType } from '@/shared/types'
import { useRouter } from 'vue-router'

export function useAuthProvider() {
	const store = useStore();
	const router = useRouter()

	const user = computed<IUser>(() => store.getters['user/getUser'])
	const userFetchState = computed<StateType>(() => store.getters['user/getUserFetchState'])

	const checkUserAuth = () => {
		store.dispatch('user/getUser').then(() => {
			if (!userFetchState.value || userFetchState.value === 'ERROR') {
				router.push('/authentication')
			} else if (userFetchState.value && userFetchState.value === 'SUCCESS') {
				store.commit('user/setUserSettings', user.value, { root: true })
				if (
					!user.value.weight ||
					!user.value.height ||
					!user.value.gender ||
					!user.value.birthday ||
					!user.value.trainCount ||
					!user.value.level ||
					!user.value.purpose ||
					!user.value.placement
				) {
					router.push('/on-boarding')
				}
			}
		})
	}

	return {
		user,
		userFetchState,
		checkUserAuth,
	}
}
