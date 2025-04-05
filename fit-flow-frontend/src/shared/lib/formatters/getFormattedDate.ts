export const getFormattedDate = (date: string) => {
	const year = date.split('T')[0].split('-')[0]
	const month = date.split('T')[0].split('-')[1]
	const day = date.split('T')[0].split('-')[2]

	return `${day}.${month}.${year}`
}
