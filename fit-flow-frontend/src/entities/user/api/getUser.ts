import { http } from '@/shared/api'

export const getUser = async () => {
	try {
		return await http.get('/auth/user')
	} catch (e) {
		console.error(e)
	}
}
