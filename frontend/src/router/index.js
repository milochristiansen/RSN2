import { createRouter, createWebHistory } from 'vue-router'
import PublicHeaders from '@/views/PublicHeaders.vue'
import Login from '@/views/public/Login.vue'

const routes = [
	// Public Pages
	{
		path: '/',
		name: 'PublicHeaders',
		component: PublicHeaders,
		children:[
			{
				path: '',
				name: 'Login',
				component: Login
			},
			{
				path: '/confirm',
				name: 'ConfirmEmail',
				component: () => import('@/views/public/ConfirmEmail.vue')
			},
			{
				path: '/confirm-email',
				name: 'EmailConfirmed',
				component: () => import('@/views/public/EmailConfirmed.vue'),
				props: route => ({ token: route.query.token }),
			},
			{
				path: '/delete-email',
				name: 'DeleteEmail',
				component: () => import('@/views/public/DeleteEmail.vue'),
				props: route => ({ token: route.query.token }),
			},
			{
				path: '/forgotpass',
				name: 'ForgotPassword',
				component: () => import('@/views/public/ForgotPassword.vue')
			},
			{
				path: '/newuser',
				name: 'NewUser',
				component: () => import('@/views/public/NewUser.vue')
			}
		]
	},
	// Private Pages
	{
		path: '/user',
		name: 'PrivateHeaders',
		component: () => import('@/views/PrivateHeaders.vue'),
		children:[
			{
				path: 'settings',
				name: 'UserSettings',
				component: () => import('@/views/private/UserSettings.vue')
			},
			{
				path: 'unread',
				name: 'Unread',
				component: () => import('@/views/private/Unread.vue')
			},
			{
				path: 'feeds',
				name: 'FeedList',
				component: () => import('@/views/private/FeedList.vue')
			},
			{
				path: 'feed-details',
				name: 'FeedDetails',
				component: () => import('@/views/private/FeedDetails.vue'),
				props: route => ({ id: route.query.id }),
			},
			{
				path: 'logout',
				name: 'Logout',
				component: () => import('@/views/private/Logout.vue')
			},
		]
	}
]

const router = createRouter({
	history: createWebHistory(process.env.BASE_URL),
	routes
})

export default router
