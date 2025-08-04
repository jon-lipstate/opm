import { defineStore } from 'pinia'
import { api } from 'boot/axios'
import { Notify } from 'quasar'
import { devLog, expectAuth } from 'src/utils/utils'
import { useUserStore } from './user-store'
import { walk } from 'vue/compiler-sfc'

export const useApiStore = defineStore('api', {
	state: () => ({
		loading: {
			packages: false,
			package: false,
			tags: false,
			user: false,
			bookmarks: false,
		},
	}),

	actions: {
		// Error handler
		handleError(error, defaultMessage = 'An error occurred') {
			console.error(error)
			const message = error.response?.data?.message || error.message || defaultMessage
			Notify.create({
				type: 'negative',
				message,
				position: 'top',
			})
			throw error
		},

		// Package endpoints
		async fetchPackages(params = {}) {
			this.loading.packages = true
			try {
				const query = '/packages'
				devLog(`GET: ${query}`, params)
				const response = await api.get(query, { params })
				devLog('Fetch Packages Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch packages')
			} finally {
				this.loading.packages = false
			}
		},

		async searchPackages(query, params = {}) {
			this.loading.packages = true
			try {
				const endpoint = '/packages/search'
				const searchParams = { q: query, ...params }
				devLog(`GET: ${endpoint}`, searchParams)
				const response = await api.get(endpoint, { params: searchParams })
				devLog('Search Packages Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Search failed')
			} finally {
				this.loading.packages = false
			}
		},

		async fetchPackage(userSlug, pkgSlug) {
			this.loading.package = true
			try {
				const query = `/packages/${userSlug}/${pkgSlug}`
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch Package Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to fetch package')
			} finally {
				this.loading.package = false
			}
		},

		async fetchPackageReadme(packageId) {
			try {
				const query = '/readme'
				const params = { package_id: packageId }
				devLog(`GET: ${query}`, params)
				const response = await api.get(query, { params })
				devLog('Fetch README Response:', response.data)
				return response.data
			} catch (error) {
				console.warn('fetchPackageReadme', error.status) // we expect to 404 here sometimes, so dont show a notify
			}
		},

		async createPackage(packageData) {
			if (!expectAuth()) return

			try {
				const query = '/packages'
				devLog(`POST: ${query}`, packageData)
				const response = await api.post(query, packageData)
				devLog('Create Package Response:', response.data)
				Notify.create({
					type: 'positive',
					message: 'Package created successfully',
				})
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to create package')
			}
		},

		// Bookmark endpoints
		async bookmarkPackage(package_id) {
			if (!expectAuth()) return

			try {
				const query = `/packages/bookmark?package_id=${package_id}`
				devLog(`POST: ${query}`)
				const response = await api.post(query)
				devLog('Bookmark Package Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to bookmark package')
			}
		},

		async unbookmarkPackage(package_id) {
			if (!expectAuth()) return

			try {
				const query = `/packages/bookmark?package_id=${package_id}`
				devLog(`DELETE: ${query}`)
				const response = await api.delete(query)
				devLog('Unbookmark Package Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to remove bookmark')
			}
		},

		// Tag endpoints
		async fetchTags(params = {}) {
			this.loading.tags = true
			try {
				const query = '/tags'
				devLog(`GET: ${query}`, params)
				const response = await api.get(query, { params })
				devLog('Fetch Tags Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch tags')
			} finally {
				this.loading.tags = false
			}
		},

		async addPackageTag(packageId, tagName) {
			if (!expectAuth()) return

			try {
				const query = '/tags'
				const data = {
					package_id: packageId,
					tag_name: tagName,
				}
				devLog(`POST: ${query}`, data)
				const response = await api.post(query, data)
				devLog('Add Package Tag Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to add tag')
			}
		},

		async voteOnTag(packageId, tagId, vote) {
			if (!expectAuth()) return

			try {
				const query = '/tags/vote'
				const data = {
					package_id: packageId,
					tag_id: tagId,
					vote,
				}
				devLog(`POST: ${query}`, data)
				const response = await api.post(query, data)
				devLog('Vote On Tag Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to vote on tag')
			}
		},

		// Flag/moderation endpoints
		async flagPackage(packageId, reason, details) {
			if (!expectAuth()) return

			try {
				const query = '/flags'
				const data = { package_id: packageId, reason, details }
				devLog(`POST: ${query}`, data)
				const response = await api.post(query, data)
				devLog('Flag Package Response:', response.data)
				Notify.create({
					type: 'positive',
					message: 'Package flagged for review',
				})
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to flag package')
			}
		},

		async fetchPackageFlags(packageId) {
			try {
				const query = '/flags'
				const params = { package_id: packageId }
				devLog(`GET: ${query}`, params)
				const response = await api.get(query, { params })
				devLog('Fetch Package Flags Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch flags')
			}
		},

		async fetchFlagStats(packageId) {
			try {
				const query = '/flags/stats'
				const params = { package_id: packageId }
				devLog(`GET: ${query}`, params)
				const response = await api.get(query, { params })
				devLog('Fetch Flag Stats Response:', response.data)
				return response.data || {}
			} catch (error) {
				this.handleError(error, 'Failed to fetch flag statistics')
			}
		},

		// Auth endpoints
		async fetchCurrentUser() {
			this.loading.user = true
			try {
				const query = '/auth/me'
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch Current User Response:', response.data)
				return response.data
			} catch (error) {
				// Don't show error for 401 - user is simply not logged in
				if (error.response?.status !== 401) {
					this.handleError(error, 'Failed to fetch user')
				}
				return null
			} finally {
				this.loading.user = false
			}
		},

		async logout() {
			try {
				const query = '/auth/logout'
				devLog(`POST: ${query}`)
				await api.post(query)
				Notify.create({
					type: 'positive',
					message: 'Logged out successfully',
				})
			} catch (error) {
				this.handleError(error, 'Failed to logout')
			}
		},

		// User endpoints
		async fetchUserBookmarks() {
			if (!expectAuth()) return

			this.loading.bookmarks = true
			try {
				const query = '/users/me/bookmarks'
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch User Bookmarks Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch bookmarks')
			} finally {
				this.loading.bookmarks = false
			}
		},

		async fetchUserPackages() {
			try {
				const query = '/users/me/packages'
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch User Packages Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch your packages')
			}
		},

		async updatePackage(packageId, updates) {
			if (!expectAuth()) return

			try {
				const query = `/packages/${packageId}`
				devLog(`PUT: ${query}`, updates)
				const response = await api.put(query, updates)
				devLog('Update Package Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to update package')
			}
		},

		async deletePackage(packageId) {
			if (!expectAuth()) return

			try {
				const query = `/packages/${packageId}`
				devLog(`DELETE: ${query}`)
				const response = await api.delete(query)
				devLog('Delete Package Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to delete package')
			}
		},

		async fetchUserFlags() {
			if (!expectAuth()) return

			try {
				const query = '/users/me/flags'
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch User Flags Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch your flags')
			}
		},

		async deleteFlag(flagId) {
			if (!expectAuth()) return

			try {
				const query = `/flags/${flagId}`
				devLog(`DELETE: ${query}`)
				const response = await api.delete(query)
				devLog('Delete Flag Response:', response.data)
				Notify.create({
					type: 'positive',
					message: 'Report retracted successfully',
				})
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to retract report')
			}
		},

		async updateUserProfile(updates) {
			try {
				const query = '/users/me'
				devLog(`PUT: ${query}`, updates)
				const response = await api.put(query, updates)
				devLog('Update Profile Response:', response.data)

				// Fetch the updated user data
				const userStore = useUserStore()
				await userStore.fetchUser()

				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to update profile')
			}
		},

		async checkSlugAvailability(slug) {
			try {
				const query = '/users/check-user-slug'
				devLog(`GET: ${query}?slug=${slug}`)
				const response = await api.get(query, { params: { slug } })
				devLog('Check slug Response:', response.data)
				return response.data
			} catch (error) {
				console.error('Failed to check slug:', error)
				throw error
			}
		},

		async fetchRepositoryMetadata(repoUrl) {
			if (!expectAuth()) return

			try {
				const query = '/repository/metadata'
				devLog(`GET: ${query}?url=${repoUrl}`)
				const response = await api.get(query, { params: { url: repoUrl } })
				devLog('Repository Metadata Response:', response.data)
				return response.data
			} catch (error) {
				// Don't show error notification - just return null
				console.error('Failed to fetch repository metadata:', error)
				return null
			}
		},

		// Moderator flag endpoints
		async fetchAllFlags(status = 'pending') {
			if (!expectAuth()) return

			try {
				const query = `/flags/all?status=${status}`
				devLog(`GET: ${query}`)
				const response = await api.get(query)
				devLog('Fetch All Flags Response:', response.data)
				return response.data || []
			} catch (error) {
				this.handleError(error, 'Failed to fetch flags')
			}
		},

		async resolveFlag(flagId, status) {
			if (!expectAuth()) return

			try {
				const query = `/flags/${flagId}/resolve`
				const data = { status }
				devLog(`PUT: ${query}`, data)
				const response = await api.put(query, data)
				devLog('Resolve Flag Response:', response.data)
				return response.data
			} catch (error) {
				this.handleError(error, 'Failed to resolve flag')
			}
		},
	},
})
