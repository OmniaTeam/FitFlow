import { useStore } from 'vuex'
import { computed } from 'vue'

export function useLoginWithGoogle() {
	const store = useStore()

	const userLoginState =  computed(() => store.getters['user/getUserLoginState'])

	const loginHandler = () => {
		store.dispatch('user/loginWithGoogle')
	}

	return {
		userLoginState,
		loginHandler,
	}
}
