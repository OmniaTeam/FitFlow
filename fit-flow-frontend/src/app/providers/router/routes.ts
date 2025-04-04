export const routes = [
    {
        path: '/',
        name: 'main-page',
        component: () => import('@/pages/main-page/main-page.vue'),
    },
    {
        path: '/auth',
        name: 'auth-page',
        component: () => import('@/pages/auth-page/auth-page.vue'),
    },
]
