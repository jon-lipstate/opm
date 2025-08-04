import { defineBoot } from '#q-app/wrappers'
import axios from 'axios'
import { useUserStore } from 'stores/user-store'

// Set base URL based on environment
const apiUrl = import.meta.env.VITE_API_URL
const api = axios.create({
	baseURL: apiUrl,
	withCredentials: true, // Important for cookies
})

export default defineBoot(({ app, router, store }) => {
	// for use inside Vue files (Options API) through this.$axios and this.$api
	app.config.globalProperties.$axios = axios
	app.config.globalProperties.$api = api

	// Request interceptor
	api.interceptors.request.use(
		(config) => {
			// Token is handled by httpOnly cookies, no need to add Authorization header
			return config
		},
		(error) => {
			return Promise.reject(error)
		},
	)

	// Response interceptor
	api.interceptors.response.use(
		(response) => {
			// Reset offline status if we get a successful response
			const userStore = useUserStore(store)
			if (userStore.isOffline) {
				userStore.setOffline(false)
			}
			return response
		},
		async (error) => {
			const originalRequest = error.config
			const userStore = useUserStore(store)

			// Handle network errors
			if (!error.response) {
				userStore.setOffline(true)
				return Promise.reject(error)
			}

			// Handle 401 Unauthorized
			if (error.response.status === 401 && !originalRequest._retry) {
				// Skip auth redirect for login-related routes
				const publicRoutes = ['/login', '/auth/callback', '/']
				if (publicRoutes.includes(router.currentRoute.value.path)) {
					return Promise.reject(error)
				}

				originalRequest._retry = true

				// Clear user state
				userStore.logout()

				// Save current route for redirect after login
				const currentPath = router.currentRoute.value.fullPath
				if (!publicRoutes.includes(currentPath)) {
					sessionStorage.setItem('redirectAfterLogin', currentPath)
				}

				// Redirect to login
				router.push('/login')
				return Promise.reject(error)
			}

			return Promise.reject(error)
		},
	)
})

export { api }
