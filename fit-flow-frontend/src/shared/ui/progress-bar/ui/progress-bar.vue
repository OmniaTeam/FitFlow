<template>
	<div class="progress-bar-container">
		<div class="back-button" @click="onBackClick">
			<div class="arrow-icon">←</div>
		</div>
		<Progress :percentage="progressPercentage" :stroke-width="4" :show-pivot="false" />
	</div>
</template>

<script lang="ts" setup>
import { Progress } from 'vant';
import { computed } from 'vue';

const props = defineProps<{
	activeQuestion: number;
	totalQuestions: number;
}>();

const emit = defineEmits<{
	(e: 'back'): void;
}>();

const progressPercentage = computed(() => {
	if (props.totalQuestions <= 1) return 0;

	return Math.min(100, Math.floor((props.activeQuestion / (props.totalQuestions - 1)) * 100));
});

const onBackClick = () => {
	emit('back');
};
</script>

<style lang="scss" scoped>
.progress-bar-container {
	display: flex;
	align-items: center;
	width: 100%;
	gap: 10px;
}

.back-button {
	cursor: pointer;
	padding: 8px;
}

.arrow-icon {
	color: #E6E6E6;
	font-size: 20px;
}

:deep(.van-progress) {
	flex-grow: 1;
	--van-progress-background: #000000;
	--van-progress-color: #E6E6E6;
}
</style>
