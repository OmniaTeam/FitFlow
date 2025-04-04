import type { IUser } from '@/entities/user/model/interfaces'
import type { StateType } from '@/shared/types'
import type { Commit } from 'vuex'
import { getUser, loginWithGoogle } from '@/entities/user/api'

type UserModuleState = {
	user: IUser | null,
	userSettings: IUser | null,
	userFetchState: StateType,
	userLoginState: StateType,
}

export const module = {
	namespaced: true,
	state: {
		user: null,
		userSettings: {
			id: null,
			firstName: null,
			lastName: null,
		},
		userFetchState: null,
		userLoginState: null,
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
	},
	mutations: {
		setUserSettings(state: UserModuleState, data: Partial<IUser> | IUser | null) {
			if (!data) {
				state.userSettings = {
					id: null,
					firstName: null,
					lastName: null,
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
	},
	actions: {
		async loginWithGoogle({ commit }: { commit: Commit }) {
			try {
				commit('setUserLoginState', 'PENDING')

				const result = await loginWithGoogle()

				if (result && result.status !== 200 && result.status !== 301) {
					commit('setUserLoginState', 'ERROR')
				} else {
					commit('setUserLoginState', 'SUCCESS')
				}

				return result
			} catch (e) {
				commit('setUserLoginState', 'ERROR')
				console.error(e)
			}
		},
		async getUser({ commit }: { commit: Commit }) {
			try {
				commit('setUserFetchState', 'PENDING')

				const result = await getUser()

				if (result && result.status === 200) {
					commit('setUserFetchState', 'SUCCESS')
				} else {
					commit('setUserLoginState', 'ERROR')
				}

				return result
			} catch (e) {
				commit('setUserFetchState', 'ERROR')
				console.error(e)
			}
		},
	},
}
