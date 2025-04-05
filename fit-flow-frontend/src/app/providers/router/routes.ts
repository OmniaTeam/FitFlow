export const routes = [
    {
        path: '/',
        name: 'home',
        component: () => import('@/pages/main-page/main-page.vue'),
    },
    {
        path: '/authentication',
        name: 'auth',
        component: () => import('@/pages/auth-page/auth-page.vue'),
    },
	{
		path: '/on-boarding',
		name: 'on-boarding',
		component: () => import('@/pages/on-boarding-page/on-boarding-page.vue'),
	},
	{
		path: '/train',
		name: 'train',
		component: () => import('@/pages/train-page/train-page.vue'),
	},
	{
		path: '/food',
		name: 'food',
		component: () => import('@/pages/food-page/food-page.vue'),
	},
	{
		path: '/chat',
		name: 'chat',
		component: () => import('@/pages/chat-page/chat-page.vue'),
	},
	{
		path: '/account',
		name: 'account',
		component: () => import('@/pages/account-page/account-page.vue'),
	}
]
