<template>
	<q-page class="flex flex-center">
		<q-card class="q-pa-md" style="min-width: 350px">
			<q-card-section>
				<div class="text-h6 text-center">Sign In to OPM</div>
				<div class="text-subtitle2 text-center text-grey q-mt-sm">Odin Package Registry</div>
			</q-card-section>

			<q-card-section>
				<div class="q-gutter-md">
					<q-btn
						label="Sign in with GitHub"
						icon="fab fa-github"
						color="grey-9"
						class="full-width"
						size="lg"
						@click="loginWithGitHub"
						:loading="loading"
					/>

					<q-btn
						label="Sign in with Discord"
						icon="fab fa-discord"
						color="indigo"
						class="full-width"
						size="lg"
						@click="loginWithDiscord"
						:loading="loading"
						:disable="true"
					/>

					<div class="text-caption text-center text-grey">Discord login coming soon</div>
				</div>
			</q-card-section>
		</q-card>
	</q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from 'stores/user-store'

const userStore = useUserStore()
const loading = ref(false)

const loginWithGitHub = () => {
	loading.value = true
	userStore.login('github')
}

const loginWithDiscord = () => {
	loading.value = true
	userStore.login('discord')
}
</script>
