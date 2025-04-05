import { EGender, ELevel, EPlacement, EPurpose } from '../enums'

export interface IUser {
	id: number | null,
	firstName: string | null,
	lastName: string | null,
	weight: number | null,
	height: number | null,
	gender: EGender | null,
	birthday: string | null,
	purpose: EPurpose | null,
	placement: EPlacement | null,
	level: ELevel | null,
	trainCount: number | null,
	foodPrompt: string | null,
}
