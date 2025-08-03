const routes = [
	{
		path: '/',
		component: () => import('layouts/MainLayout.vue'),
		children: [
			{ path: '', component: () => import('pages/IndexPage.vue') },
			{ path: 'login', component: () => import('pages/LoginPage.vue') },
			{
				path: 'packages/:userAlias/:packageSlug',
				component: () => import('pages/PackageDetailPage.vue'),
			},
			{ path: 'submit', component: () => import('pages/SubmitPackagePage.vue') },
			{ path: 'my-packages', component: () => import('pages/MyPackagesPage.vue') },
			{ path: 'profile', component: () => import('pages/ProfilePage.vue') },
		],
	},

	// Always leave this as last one,
	// but you can also remove it
	{
		path: '/:catchAll(.*)*',
		component: () => import('pages/ErrorNotFound.vue'),
	},
]

export default routes
