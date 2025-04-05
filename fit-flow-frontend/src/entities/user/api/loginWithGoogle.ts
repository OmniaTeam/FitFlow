import { http } from '@/shared/api'

export const loginWithGoogle = async () => {
	try {
		return await http.get('/auth/google/login')
	} catch (e) {
		console.error(e)
	}
}
