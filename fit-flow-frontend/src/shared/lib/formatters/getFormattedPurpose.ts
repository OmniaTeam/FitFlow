export const getFormattedPurpose = (purpose: string) => {
	switch (purpose) {
		case 'gain':
			return 'Набор мышечной массы'
		case 'slim':
			return 'Похудение'
		case 'stable':
			return 'Поддержание формы'
	}
}
