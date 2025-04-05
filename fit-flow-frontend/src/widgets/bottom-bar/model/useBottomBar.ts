import { computed } from 'vue'
import { useRouter } from 'vue-router'

export function useBottomBar() {
	const router = useRouter()

	console.log(router)

	const activeTab = computed(() => router.currentRoute.value.name)

	const changeTab = async (path: string) => {
		await router.push(path)
	}

	return {
		activeTab,
		changeTab,
	}
}
