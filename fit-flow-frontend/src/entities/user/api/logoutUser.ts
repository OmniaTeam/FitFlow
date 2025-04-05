import { http } from '@/shared/api'

export const logoutUser = async () => {
	try {
		return await http.get('/auth/logout')
	} catch (e) {
		console.error(e);
	}
}
