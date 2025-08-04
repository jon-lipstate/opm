<template>
	<router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from 'stores/user-store'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

onMounted(async () => {
	// Check if we're returning from OAuth
	const handled = await userStore.handleAuthCallback(router)

	// Only try to fetch user if:
	// 1. Not handling OAuth callback
	// 2. We have a stored user (indicating previous session)
	// 3. Not already authenticated in current session
	const hasStoredUser = localStorage.getItem('user') !== null
	if (!handled && hasStoredUser && !userStore.isAuthenticated) {
		await userStore.fetchUser()
	}
})
</script>
