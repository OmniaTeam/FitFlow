import { computed, type EmitFn } from 'vue'
import { Account, Chat, Food, Home, Train } from '@/shared/ui/icons'

export function useNavTab(icon: string, emit: EmitFn) {
	const iconComponent = computed(() => {
		switch (icon) {
			case 'home': return Home
			case 'train': return Train
			case 'food': return Food
			case 'chat': return Chat
			case 'account': return Account
			default: return null
		}
	})

	const onClick = () => {
		emit('click')
	}

	return {
		iconComponent,
		onClick,
	}
}
