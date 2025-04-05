export const getFormattedLevel = (level: string) => {
	switch (level) {
		case 'pro':
			return 'Профессионал'
		case 'medium':
			return 'Средний уровень'
		case 'new':
			return 'Новичок'
	}
}
