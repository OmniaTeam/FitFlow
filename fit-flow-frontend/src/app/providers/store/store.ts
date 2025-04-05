import { userModule } from '@/entities/user'
import { createStore } from 'vuex'

export const store = createStore({
    modules: {
		['user']: userModule.module,
	},
})
