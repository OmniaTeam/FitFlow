import type { IUser } from '@/entities/user/model/interfaces'
import type { StateType } from '@/shared/types'
import type { Commit } from 'vuex'
import { editUser, getUser, loginWithGoogle, loginWithVkid, logoutUser } from '@/entities/user/api'
import type { EditUserInput } from '@/entities/user'

type UserModuleState = {
	user: IUser | null,
	userSettings: IUser | null,
	userFetchState: StateType,
	userLoginState: StateType,
	userEditState: StateType,
	userLogoutState: StateType,
}

export const module = {
	namespaced: true,
	state: {
		user: null,
		userSettings: {
			id: null,
			firstName: null,
			lastName: null,
			weight: null,
			height: null,
			gender: null,
			birthday: null,
			purpose: null,
			placement: null,
			level: null,
			trainCount: null,
			foodPrompt: null,
		},
		userFetchState: null,
		userLoginState: null,
		userEditState: null,
		userLogoutState: null,
	} as UserModuleState,
	getters: {
		getUser(state: UserModuleState) {
			return state.user
		},
		getUserSettings(state: UserModuleState) {
			return state.userSettings
		},
		getUserFetchState(state: UserModuleState) {
			return state.userFetchState
		},
		getUserEditState(state: UserModuleState) {
			return state.userEditState
		},
		getUserLogoutState(state: UserModuleState) {
			return state.userLogoutState
		}
	},
	mutations: {
		setUser(state: UserModuleState, data: IUser) {
			state.user = data
		},
		setUserSettings(state: UserModuleState, data: Partial<IUser> | IUser | null) {
			if (!data) {
				state.userSettings = {
					id: null,
					firstName: null,
					lastName: null,
					weight: null,
					height: null,
					gender: null,
					birthday: null,
					purpose: null,
					placement: null,
					level: null,
					trainCount: null,
					foodPrompt: null,
				}
			} else {
				Object.assign(state.userSettings as IUser, data)
			}
		},
		setUserFetchState(state: UserModuleState, data: StateType) {
			state.userFetchState = data
		},
		setUserLoginState(state: UserModuleState, data: StateType) {
			state.userLoginState = data
		},
		setUserEditState(state: UserModuleState, data: StateType) {
			state.userEditState = data
		},
		setUserLogoutState(state: UserModuleState, data: StateType) {
			state.userLogoutState = data
		}
	},
	actions: {
		async loginWithGoogle({ commit }: { commit: Commit }) {
			try {
				commit('setUserLoginState', 'PENDING')

				const result = await loginWithGoogle()

				if (result && result.status !== 200 && result.status !== 301) {
					commit('setUserLoginState', 'ERROR')
				} else if (result && result.status === 200) {
					commit('setUserLoginState', 'SUCCESS')
					window.location.href = result.data.redirect_url
				}

				return result
			} catch (e) {
				commit('setUserLoginState', 'ERROR')
				console.error(e)
			}
		},
		async loginWithVkid({ commit }: { commit: Commit }) {
			try {
				commit('setUserLoginState', 'PENDING')

				const result = await loginWithVkid()

				if (result && result.status !== 200 && result.status !== 301) {
					commit('setUserLoginState', 'ERROR')
				} else if (result && result.status === 200) {
					commit('setUserLoginState', 'SUCCESS')
					window.location.href = result.data.redirect_url
				}

				return result
			} catch (e) {
				commit('setUserLoginState', 'ERROR')
				console.error(e)
			}
		},
		async logout({ commit }: { commit: Commit }) {
			try {
				commit('setUserLogoutState', 'PENDING')

				const result = await logoutUser()

				if (result && result.status !== 200) {
					commit('setUserLogoutState', 'ERROR')
				} else if (result && result.status === 200) {
					commit('setUserLoginState', 'SUCCESS')
					window.location.reload()
				}

				return result
			} catch (e) {
				commit('setUserLogoutState', 'ERROR')
			}
		},
		async getUser({ commit }: { commit: Commit }) {
			try {
				commit('setUserFetchState', 'PENDING')

				const result = await getUser()

				if (result && result.status === 200) {
					commit('setUser', {
						id: result.data.id,
						firstName: result.data.first_name,
						lastName: result.data.last_name,
						weight: result.data.weight,
						height: result.data.height,
						gender: result.data.gender,
						birthday: result.data.birthday,
						purpose: result.data.purpose,
						placement: result.data.placement,
						level: result.data.level,
						trainCount: result.data.training_count,
						foodPrompt: result.data.food_prompt,
					})
					commit('setUserFetchState', 'SUCCESS')

				} else {
					commit('setUserFetchState', 'ERROR')
				}

				return result
			} catch (e) {
				commit('setUserFetchState', 'ERROR')
				console.error(e)
			}
		},
		async editUser({ commit }: { commit: Commit }, input: EditUserInput) {
			try {
				commit('setUserEditState', 'PENDING')
				const result = await editUser(input)

				if (result && result.status === 200) {
					console.log(result.data, 'for success')
					commit('setUserEditState', 'SUCCESS')
				} else {
					console.log(result, 'for error')
					commit('setUserEditState', 'ERROR')
				}

				return result
			} catch (e) {
				console.error(e)
				commit('setUserEditState', 'ERROR')
			}
		}
	},
}
