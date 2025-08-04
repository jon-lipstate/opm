import { defineStore } from 'pinia'
import { useApiStore } from './api-store'
import { Notify } from 'quasar'

export const useUserStore = defineStore('user', {
	state: () => ({
		user: JSON.parse(localStorage.getItem('user')) || null,
		isAuthenticated: !!localStorage.getItem('user'),
		isOffline: false,
		isLoading: false,
	}),

	getters: {
		username: (state) => state.user?.username || '',
		displayName: (state) => state.user?.display_name || state.user?.username || '',
		avatarUrl: (state) => state.user?.avatar_url || '',
		alias: (state) => state.user?.alias || state.user?.username || '',
		isLoggedIn: (state) => state.isAuthenticated && state.user !== null,
		isModerator: (state) => state.user?.is_moderator || false,
	},

	actions: {
		async fetchUser() {
			try {
				this.isLoading = true
				const apiStore = useApiStore()
				const user = await apiStore.fetchCurrentUser()
				if (user) {
					this.setAuthState(user)
				}
				return this.user
			} catch (_error) {
				this.clearAuthState()
				return null
			} finally {
				this.isLoading = false
			}
		},

		async login(provider) {
			const apiURL = import.meta.env.VITE_API_URL
			// OAuth login is handled by redirecting to the server
			window.location.href = `${apiURL}/auth/${provider}`
		},

		async logout() {
			try {
				const apiStore = useApiStore()
				await apiStore.logout()
			} catch (error) {
				console.error('Logout error:', error)
			} finally {
				// Clear auth state
				this.clearAuthState()
			}
		},

		setOffline(offline) {
			this.isOffline = offline
			if (offline) {
				Notify.create({
					type: 'negative',
					message: 'Network connection lost',
					position: 'top',
					timeout: 0,
					actions: [{ label: 'Dismiss', color: 'white' }],
				})
			}
		},

		setAuthState(user) {
			this.user = user
			this.isAuthenticated = true

			// Save to localStorage
			localStorage.setItem('user', JSON.stringify(user))
		},

		clearAuthState() {
			this.user = null
			this.isAuthenticated = false

			// Clear from localStorage
			localStorage.removeItem('user')

			// Clear any stored redirect
			sessionStorage.removeItem('redirectAfterLogin')
		},

		// Handle OAuth callback
		async handleAuthCallback(router) {
			const urlParams = new URLSearchParams(window.location.search)
			if (urlParams.get('auth') === 'success') {
				// Clean up URL
				window.history.replaceState({}, document.title, window.location.pathname)

				// Fetch user data
				await this.fetchUser()

				// Check for redirect
				const redirect = sessionStorage.getItem('redirectAfterLogin')
				if (redirect) {
					sessionStorage.removeItem('redirectAfterLogin')
					router.push(redirect)
				}

				// Show success message
				Notify.create({
					type: 'positive',
					message: 'Successfully logged in!',
					position: 'top',
				})

				return true
			}
			return false
		},
	},
})
