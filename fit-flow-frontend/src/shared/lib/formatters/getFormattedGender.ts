export const getFormattedGender = (gender: string) => {
	switch (gender) {
		case 'male':
			return 'Мужчина '
		case 'female':
			return 'Женщина'
	}
}
