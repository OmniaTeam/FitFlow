import { useStore } from 'vuex'
import { computed } from 'vue'

export function useLoginWithVkid() {
	const store = useStore()

	const userLoginState =  computed(() => store.getters['user/getUserLoginState'])

	const loginHandler = () => {
		store.dispatch('user/loginWithVkid')
	}

	return {
		userLoginState,
		loginHandler,
	}
}
