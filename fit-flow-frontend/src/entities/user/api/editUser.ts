import { http } from '@/shared/api'
import type { EditUserInput } from '../model'

export const editUser = async (input: EditUserInput) => {
	try {
		return await http.put('/gym/users', {
			"birthday": input.birthday,
			"first_name": input.firstName,
			"food_prompt": input.foodPrompt,
			"gender": input.gender,
			"height": Number(input.height),
			"last_name": input.lastName,
			"level": input.level,
			"placement": input.placement,
			"purpose": input.purpose,
			"training_count": Number(input.trainingCount),
			"weight": Number(input.weight)
		})
	} catch (e) {
		console.error(e)
	}
}
