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

	// If not OAuth callback and not already authenticated, try to fetch current user
	if (!handled && !userStore.isAuthenticated) {
		await userStore.fetchUser()
	}
})
</script>
