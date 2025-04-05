<template>
	<svg
		width="25"
		height="24"
		viewBox="0 0 25 24"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
	>
		<path
			d="M22 17H21.5V16C21.4978 11.418 18.0536 7.56911 13.5 7.06V5.5C13.5 5.22386 13.2761 5 13 5H12C11.7239 5 11.5 5.22386 11.5 5.5V7.06C6.94641 7.56911 3.50217 11.418 3.5 16V17H3C2.72386 17 2.5 17.2239 2.5 17.5V18.5C2.5 18.7761 2.72386 19 3 19H22C22.2761 19 22.5 18.7761 22.5 18.5V17.5C22.5 17.2239 22.2761 17 22 17Z"
			:fill="active ? '#010101' : '#E6E6E6'"
		/>
	</svg>
</template>

<script lang="ts" setup>
defineProps({
	active: {
		type: Boolean,
		required: true,
	}
})
</script>
