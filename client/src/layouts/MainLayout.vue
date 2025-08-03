<template>
	<q-layout view="lHh Lpr lFf">
		<q-header elevated :class="$q.dark.isActive ? 'bg-odin-dark' : 'bg-odin-primary'">
			<q-toolbar>
				<q-toolbar-title>
					<div class="row items-center no-wrap text-white">
						<router-link to="/">
							<img
								src="/odin-logo.svg"
								alt="Odin"
								style="height: 32px; width: auto"
								class="q-mt-xs q-mr-md"
							/>
						</router-link>
						<span class="text-subtitle1 text-grey-5">Package Registry</span>
					</div>

					<div></div>
				</q-toolbar-title>

				<q-btn-dropdown v-if="userStore.isLoggedIn" flat no-caps>
					<template v-slot:label>
						<div class="row items-center no-wrap">
							<q-avatar size="32px" class="q-mr-sm">
								<img v-if="userStore.avatarUrl" :src="userStore.avatarUrl" />
								<q-icon v-else name="account_circle" />
							</q-avatar>
							{{ userStore.displayName }}
						</div>
					</template>

					<q-list>
						<q-item clickable v-close-popup to="/profile">
							<q-item-section avatar>
								<q-icon name="person" />
							</q-item-section>
							<q-item-section>Profile</q-item-section>
						</q-item>

						<q-item clickable v-close-popup to="/my-packages">
							<q-item-section avatar>
								<q-icon name="inventory_2" />
							</q-item-section>
							<q-item-section>My Packages</q-item-section>
						</q-item>

						<q-item clickable v-close-popup to="/submit">
							<q-item-section avatar>
								<q-icon name="add" />
							</q-item-section>
							<q-item-section>Submit Package</q-item-section>
						</q-item>

						<q-separator />

						<q-item clickable v-close-popup @click="logout">
							<q-item-section avatar>
								<q-icon name="logout" />
							</q-item-section>
							<q-item-section>Logout</q-item-section>
						</q-item>
					</q-list>
				</q-btn-dropdown>

				<q-btn v-else flat label="Sign In" icon="login" to="/login" />

				<!-- Moderation Button (for moderators) -->
				<q-btn
					v-if="userStore.isModerator"
					flat
					round
					icon="shield"
					@click="rightDrawerOpen = !rightDrawerOpen"
					class="q-ml-sm"
				>
					<q-tooltip>Moderation Queue</q-tooltip>
					<q-badge v-if="pendingReportsCount > 0" color="red" floating rounded>
						{{ pendingReportsCount }}
					</q-badge>
				</q-btn>
			</q-toolbar>
		</q-header>

		<q-drawer
			v-if="userStore.isModerator"
			v-model="rightDrawerOpen"
			side="right"
			overlay
			bordered
			:width="400"
			class="bg-grey-1"
			:class="$q.dark.isActive ? 'bg-dark' : 'bg-grey-1'"
		>
			<ModerationPanel />
		</q-drawer>

		<q-page-container>
			<router-view />
		</q-page-container>

		<!-- Floating Dark Mode Toggle -->
		<q-page-sticky position="bottom-right" :offset="[18, 18]">
			<q-btn
				fab
				flat
				:icon="$q.dark.isActive ? 'dark_mode' : 'light_mode'"
				:color="$q.dark.isActive ? 'amber' : 'grey-8'"
				@click="toggleDarkMode"
			>
			</q-btn>
		</q-page-sticky>
	</q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from 'stores/user-store'
import { useRouter } from 'vue-router'
import { Dark } from 'quasar'
import ModerationPanel from 'components/ModerationPanel.vue'

const userStore = useUserStore()
const router = useRouter()

const rightDrawerOpen = ref(false)

// Mock pending reports count - in production this would come from an API
const pendingReportsCount = computed(() => {
	// This would be fetched from the server
	return userStore.isModerator ? 3 : 0
})

const logout = async () => {
	await userStore.logout()
	router.push('/')
}

const toggleDarkMode = () => {
	Dark.toggle()
	// Save preference to localStorage
	localStorage.setItem('darkMode', Dark.isActive.toString())
}
</script>
