import { http } from '@/shared/api'

export const getUser = async () => {
	try {
		return await http.get('/gym/users/profile')
	} catch (e) {
		console.error(e)
	}
}
