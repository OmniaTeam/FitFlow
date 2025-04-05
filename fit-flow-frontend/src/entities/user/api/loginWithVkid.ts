import { http } from '@/shared/api'

export const loginWithVkid = async () => {
	try {
		return await http.get('/auth/vkid/login')
	} catch (e) {
		console.error(e)
	}
}
