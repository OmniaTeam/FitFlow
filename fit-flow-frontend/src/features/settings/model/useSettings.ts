import { useStore } from 'vuex'
import { computed, ref } from 'vue'
import type { IUser } from '@/entities/user'

export function useSettings() {
	const store = useStore();

	const userSettings = computed<IUser>(() => store.getters['user/getUserSettings'])

	const showBirthdayPicker = ref(false)
	const showGenderPicker = ref(false)
	const showWeightPicker = ref(false)
	const showHeightPicker = ref(false)
	const showPurposePicker = ref(false)
	const showPlacementPicker = ref(false)
	const showLevelPicker = ref(false)
	const showTrainCountPicker = ref(false)

	const toggleBirthdayPicker = () => showBirthdayPicker.value = !showBirthdayPicker.value
	const toggleGenderPicker = () => showGenderPicker.value = !showGenderPicker.value
	const toggleWeightPicker = () => showWeightPicker.value = !showWeightPicker.value
	const toggleHeightPicker = () => showHeightPicker.value = !showHeightPicker.value
	const togglePurposePicker = () => showPurposePicker.value = !showPurposePicker.value
	const togglePlacementPicker = () => showPlacementPicker.value = !showPlacementPicker.value
	const toggleLevelPicker = () => showLevelPicker.value = !showLevelPicker.value
	const toggleTrainCountPicker = () => showTrainCountPicker.value = !showTrainCountPicker.value

	const handleNameChange = (value: string) => {
		store.commit('user/setUserSettings', { firstName: value })
	}

	const handleBirthdayChange = (values: any) => {
		const selectedDate = `${values.selectedValues[0]}-${values.selectedValues[1]}-${values.selectedValues[2]}T00:00:00Z`
		store.commit('user/setUserSettings', { birthday: selectedDate })
		showBirthdayPicker.value = false
	}

	const handleGenderChange = (values: any) => {
		const selectedValue = values.selectedValues[0]
		store.commit('user/setUserSettings', { gender: selectedValue })
		showGenderPicker.value = false
	}

	const handleWeightChange = (value: string) => {
		store.commit('user/setUserSettings', { weight: value })
	}

	const handleHeightChange = (value: string) => {
		store.commit('user/setUserSettings', { height: value })
	}

	const handlePurposeChange = (values: any) => {
		const selectedValue = values.selectedValues[0]
		store.commit('user/setUserSettings', { purpose: selectedValue })
		showPurposePicker.value = false
	}

	const handlePlacementChange = (values: any) => {
		const selectedValue = values.selectedValues[0]
		store.commit('user/setUserSettings', { placement: selectedValue })
		showPlacementPicker.value = false
	}

	const handleLevelChange = (values: any) => {
		const selectedValue = values.selectedValues[0]
		store.commit('user/setUserSettings', { level: selectedValue })
		showLevelPicker.value = false
	}

	const handleTrainCountChange = (values: any) => {
		const selectedValue = values.selectedValues[0]
		store.commit('user/setUserSettings', { trainCount: selectedValue })
		showTrainCountPicker.value = false
	}

	const handleFoodPromptChange = (value: string) => {
		store.commit('user/setUserSettings', { foodPrompt: value })
	}

	return {
		userSettings,

		showBirthdayPicker,
		showGenderPicker,
		showWeightPicker,
		showHeightPicker,
		showPurposePicker,
		showPlacementPicker,
		showLevelPicker,
		showTrainCountPicker,

		toggleBirthdayPicker,
		toggleGenderPicker,
		toggleWeightPicker,
		toggleHeightPicker,
		togglePurposePicker,
		toggleLevelPicker,
		togglePlacementPicker,
		toggleTrainCountPicker,

		handleBirthdayChange,
		handleGenderChange,
		handleWeightChange,
		handleHeightChange,
		handleNameChange,
		handlePurposeChange,
		handlePlacementChange,
		handleLevelChange,
		handleTrainCountChange,
		handleFoodPromptChange
	}
}
