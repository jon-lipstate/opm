<template>
	<q-page class="q-pa-md">
		<div class="container" style="max-width: 800px; margin: 0 auto">
			<h4 class="q-mb-lg">My Profile</h4>

			<q-card>
				<q-card-section>
					<q-form @submit="updateProfile">
						<!-- Username (read-only) -->
						<q-input
							v-model="userStore.user.username"
							label="Username"
							outlined
							readonly
							class="q-mb-md"
							hint="Your username from GitHub or Discord (cannot be changed)"
						/>

						<!-- Alias -->
						<q-input
							v-model="form.alias"
							label="URL Alias"
							outlined
							class="q-mb-md"
							:rules="aliasRules"
							hint="Used in package URLs. Must be unique and URL-safe (lowercase, letters, numbers, hyphens, underscores)"
							@update:model-value="(val) => (form.alias = val?.toLowerCase())"
						>
							<template v-slot:prepend>
								<span class="text-caption text-grey-6">pkg-odin.org/packages/</span>
							</template>
						</q-input>

						<!-- Display Name -->
						<q-input
							v-model="form.display_name"
							label="Display Name"
							outlined
							class="q-mb-md"
							:rules="displayNameRules"
							hint="Your public display name"
						/>

						<!-- Avatar URL -->
						<q-input
							v-model="form.avatar_url"
							label="Avatar URL"
							outlined
							class="q-mb-md"
							:rules="avatarRules"
							hint="URL to your avatar image (leave empty to use default)"
						>
							<template v-slot:append v-if="form.avatar_url">
								<q-avatar size="32px">
									<img :src="form.avatar_url" @error="avatarError = true" />
								</q-avatar>
							</template>
						</q-input>

						<!-- Account Info -->
						<div class="q-mt-lg q-mb-md">
							<div class="text-h6 q-mb-sm">Account Information</div>
							<q-list dense>
								<q-item>
									<q-item-section>
										<q-item-label>Member Since</q-item-label>
									</q-item-section>
									<q-item-section side>
										<q-item-label caption>{{ formatDate(userStore.user.created_at) }}</q-item-label>
									</q-item-section>
								</q-item>

								<q-item>
									<q-item-section>
										<q-item-label>Reputation</q-item-label>
									</q-item-section>
									<q-item-section side>
										<q-item-label caption>
											{{ userStore.user.reputation }} points ({{ userStore.user.reputation_rank }})
										</q-item-label>
									</q-item-section>
								</q-item>

								<q-item>
									<q-item-section>
										<q-item-label>Verified Accounts</q-item-label>
									</q-item-section>
									<q-item-section side>
										<div class="row q-gutter-xs">
											<q-icon
												v-if="userStore.user.github_verified"
												name="mdi-github"
												size="20px"
												color="positive"
											>
												<q-tooltip>GitHub Verified</q-tooltip>
											</q-icon>
											<q-icon
												v-if="userStore.user.discord_verified"
												name="mdi-discord"
												size="20px"
												color="positive"
											>
												<q-tooltip>Discord Verified</q-tooltip>
											</q-icon>
										</div>
									</q-item-section>
								</q-item>

								<q-item v-if="userStore.user.is_moderator">
									<q-item-section>
										<q-item-label>Role</q-item-label>
									</q-item-section>
									<q-item-section side>
										<q-badge color="primary">Moderator</q-badge>
									</q-item-section>
								</q-item>
							</q-list>
						</div>

						<!-- Submit Button -->
						<div class="q-mt-lg">
							<q-btn
								type="submit"
								color="primary"
								label="Update Profile"
								:loading="loading"
								:disable="!hasChanges"
							/>
						</div>
					</q-form>
				</q-card-section>
			</q-card>
		</div>
	</q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useUserStore } from 'src/stores/user-store'
import { useApiStore } from 'src/stores/api-store'

const $q = useQuasar()
const userStore = useUserStore()
const apiStore = useApiStore()

// Form data
const form = ref({
	alias: '',
	display_name: '',
	avatar_url: '',
})

const originalForm = ref({})
const loading = ref(false)
const avatarError = ref(false)

// Validation rules
const aliasRules = [
	(val) => !!val || 'Alias is required',
	(val) => val.length >= 3 || 'Must be at least 3 characters',
	(val) => val.length <= 50 || 'Must be less than 50 characters',
	(val) => /^[a-z0-9][a-z0-9-_]*[a-z0-9]$/.test(val) || 'Must start and end with a letter or number. Can contain hyphens and underscores.',
	(val) => val === val.toLowerCase() || 'Must be lowercase',
	async (val) => {
		if (val === originalForm.value.alias) return true // No need to check if unchanged
		return await checkAliasAvailability(val)
	},
]

const displayNameRules = [(val) => !val || val.length <= 255 || 'Must be less than 255 characters']

const avatarRules = [
	(val) => !val || val.startsWith('http') || 'Must be a valid URL',
	(val) => !val || !avatarError.value || 'Invalid image URL',
]

// Computed
const hasChanges = computed(() => {
	return (
		form.value.alias !== originalForm.value.alias ||
		form.value.display_name !== originalForm.value.display_name ||
		form.value.avatar_url !== originalForm.value.avatar_url
	)
})

// Methods
const formatDate = (date) => {
	return new Date(date).toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
	})
}

const checkAliasAvailability = async (alias) => {
	try {
		const response = await apiStore.checkAliasAvailability(alias)
		if (!response.available) {
			return response.reason || 'Alias is not available'
		}
		return true
	} catch (error) {
		return 'Unable to check alias availability'
	}
}

const updateProfile = async () => {
	loading.value = true
	try {
		const updates = {}

		// Only include changed fields
		if (form.value.alias !== originalForm.value.alias) {
			updates.alias = form.value.alias
		}
		if (form.value.display_name !== originalForm.value.display_name) {
			updates.display_name = form.value.display_name || undefined
		}
		if (form.value.avatar_url !== originalForm.value.avatar_url) {
			updates.avatar_url = form.value.avatar_url || undefined
		}

		await apiStore.updateUserProfile(updates)

		// Update the original form to match
		originalForm.value = { ...form.value }

		$q.notify({
			type: 'positive',
			message: 'Profile updated successfully',
		})
	} catch (error) {
		console.error('Failed to update profile:', error)
	} finally {
		loading.value = false
	}
}

// Lifecycle
onMounted(() => {
	if (userStore.user) {
		form.value = {
			alias: userStore.user.alias || '',
			display_name: userStore.user.display_name || '',
			avatar_url: userStore.user.avatar_url || '',
		}
		originalForm.value = { ...form.value }
	}
})
</script>

<style scoped>
.container {
	:deep(.q-field__prefix) {
		font-size: 0.8rem;
		opacity: 0.7;
	}
}
</style>
