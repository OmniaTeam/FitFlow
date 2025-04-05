import { useStore } from 'vuex'
import { computed } from 'vue'
import type { IUser } from '@/entities/user'
import type { StateType } from '@/shared/types'

export function useAppHeader() {
	const store = useStore()

	const user = computed<IUser>(() => store.getters['user/getUser'])
	const userFetchState = computed<StateType>(() => store.getters['user/getUserFetchState'])

	return {
		user,
		userFetchState,
	}
}
