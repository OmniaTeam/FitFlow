export const getFormattedPlacement = (placement: string) => {
	switch (placement) {
		case 'home':
			return 'Дома'
		case 'gym':
			return 'В тренажёрном зале'
	}
}
