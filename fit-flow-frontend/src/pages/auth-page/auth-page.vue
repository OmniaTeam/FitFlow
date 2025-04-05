<template>
	<div class="auth-page">
		<div class="intro-video">
			<video autoplay muted loop class="intro-video__content">
				<source src="@/shared/assets/app-video.webm" type="video/webm">
			</video>
		</div>

		<div class="auth-page__logo">
			<span class="auth-page__logo--first">Fit</span>
			<span class="auth-page__logo--second">FLOW</span>
		</div>

		<div class="auth-page__content" :class="{ 'auth-page__content--visible': showContent }">
			<div class="auth-page__main">
				<p class="auth-page__main--title">Стань лучшей версией себя</p>
				<p class="auth-page__main--subtitle">Фитнес-революция в твоем кармане! Индивидуальные планы, мотивация и
					поддержка</p>
				<div class="auth-page__main--buttons">
					<LoginWithGoogle />
					<LoginWithVkid />
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { LoginWithGoogle } from '@/features/login-with-google'
import { LoginWithVkid } from '@/features/login-with-vkid'
import { ref, onMounted } from 'vue'

const showContent = ref(false)

onMounted(() => {
	setTimeout(() => {
		showContent.value = true
	}, 200)
})
</script>

<style lang="scss" scoped>
.auth-page {
	position: relative;
	display: flex;
	flex-direction: column;
	height: calc(100vh - 128px);
	justify-content: space-between;

	.intro-video {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		z-index: 1;
		background-color: rgba(0, 0, 0, 0.7);

		&__content {
			position: absolute;
			width: 100%;
			height: 100%;
			object-fit: cover;
		}
	}

	&__logo {
		position: relative;
		z-index: 2;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 8px;
		margin-top: 30vh;
		font-family: 'SF Pro Display', sans-serif;
		font-size: 64px;
		font-weight: bold;
		color: #E6E6E6;

		&--second {
			font-style: italic;
		}
	}

	&__content {
		position: relative;
		z-index: 3;
		width: 100%;
		opacity: 0;
		transition: opacity 0.8s ease, transform 0.8s ease;

		&--visible {
			opacity: 1;
		}
	}

	&__main {
		display: flex;
		flex-direction: column;
		gap: 20px;
		align-items: center;

		&--title {
			font-size: 24px;
			text-align: center;
			color: #E6E6E6;
		}

		&--subtitle {
			font-size: 16px;
			text-align: center;
			color: #CCCCCC;
		}

		&--buttons {
			display: flex;
			flex-direction: column;
			width: 100%;
			gap: 12px
		}
	}
}
</style>
