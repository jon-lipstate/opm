<template>
	<div class="moderation-panel">
		<q-toolbar class="bg-primary text-white">
			<q-toolbar-title>Moderation Actions</q-toolbar-title>
			<q-btn flat round dense icon="refresh" @click="refreshReports" :loading="loading">
				<q-tooltip>Refresh</q-tooltip>
			</q-btn>
		</q-toolbar>

		<div class="q-pa-md">
			<!-- Filter tabs -->
			<q-tabs
				v-model="filterTab"
				dense
				class="text-grey"
				active-color="primary"
				indicator-color="primary"
				align="justify"
				narrow-indicator
			>
				<q-tab name="pending" label="Pending" />
				<q-tab name="reviewed" label="Reviewed" />
				<q-tab name="all" label="All" />
			</q-tabs>

			<q-separator class="q-mt-md q-mb-lg" />

			<!-- Reports list -->
			<div v-if="loading && reports.length === 0" class="text-center q-py-xl">
				<q-spinner size="40px" color="primary" />
			</div>

			<div v-else-if="filteredReports.length === 0" class="text-center q-py-xl text-grey-6">
				<q-icon name="check_circle" size="48px" />
				<p class="q-mt-sm">No {{ filterTab }} reports</p>
			</div>

			<q-list v-else separator>
				<q-item v-for="report in filteredReports" :key="report.id" class="q-pa-md">
					<q-item-section>
						<!-- Package info -->
						<div class="row items-center q-mb-sm">
							<router-link
								:to="`/packages/${report.package.author_slug}/${report.package.name}`"
								class="text-subtitle1 text-primary"
								target="_blank"
							>
								{{ report.package.display_name }}
							</router-link>
							<q-space />
							<q-badge
								:color="
									report.status === 'pending'
										? 'orange'
										: report.status === 'resolved'
											? 'positive'
											: 'grey'
								"
								:label="report.status"
							/>
						</div>

						<!-- Report details -->
						<div class="text-caption text-grey-6 q-mb-xs">
							Reported by {{ report.reporter.username }} • {{ formatDate(report.created_at) }}
						</div>

						<q-chip
							:color="getReasonColor(report.reason)"
							text-color="white"
							size="sm"
							class="q-mb-sm"
						>
							{{ report.reason }}
						</q-chip>

						<p v-if="report.details" class="text-body2 q-mb-md">
							{{ report.details }}
						</p>

						<!-- Action buttons -->
						<div v-if="report.status === 'pending'" class="row q-gutter-sm">
							<q-btn
								size="sm"
								color="negative"
								label="Remove Package"
								icon="delete"
								no-caps
								@click="showActionDialog('remove', report)"
							/>
							<q-btn
								size="sm"
								color="deep-orange"
								label="Ban User"
								icon="block"
								no-caps
								@click="showActionDialog('ban', report)"
							/>
							<q-btn
								size="sm"
								color="grey"
								label="Dismiss"
								icon="close"
								no-caps
								@click="showActionDialog('dismiss', report)"
							/>
						</div>

						<!-- Resolution info -->
						<div v-else class="text-caption text-grey-6 q-mt-sm">
							{{ report.status === 'resolved' ? 'Resolved' : 'Dismissed' }} by
							{{ report.resolved_by?.username }} •
							{{ formatDate(report.resolved_at) }}
						</div>
					</q-item-section>
				</q-item>
			</q-list>
		</div>

		<!-- Action confirmation dialog -->
		<q-dialog v-model="actionDialog.show">
			<q-card style="min-width: 400px">
				<q-card-section>
					<div class="text-h6">Confirm {{ actionDialog.type }}</div>
				</q-card-section>

				<q-card-section>
					<p v-if="actionDialog.type === 'remove'">
						Are you sure you want to remove the package
						<strong>{{ actionDialog.report?.package.display_name }}</strong
						>? This action cannot be undone.
					</p>
					<p v-else-if="actionDialog.type === 'ban'">
						Are you sure you want to ban the user
						<strong>{{ actionDialog.report?.package.author_username }}</strong
						>? They will not be able to submit new packages.
					</p>
					<p v-else>
						Are you sure you want to dismiss this report? The package will remain visible.
					</p>

					<q-input
						v-model="actionDialog.notes"
						type="textarea"
						label="Notes (optional)"
						outlined
						rows="3"
						class="q-mt-md"
					/>
				</q-card-section>

				<q-card-actions align="right">
					<q-btn flat label="Cancel" v-close-popup />
					<q-btn
						:color="actionDialog.type === 'dismiss' ? 'grey' : 'negative'"
						:label="actionDialog.type === 'dismiss' ? 'Dismiss Report' : 'Confirm'"
						@click="executeAction"
						:loading="actionDialog.loading"
					/>
				</q-card-actions>
			</q-card>
		</q-dialog>
	</div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Notify } from 'quasar'

const loading = ref(false)
const filterTab = ref('pending')
const reports = ref([])

const actionDialog = ref({
	show: false,
	type: null,
	report: null,
	notes: '',
	loading: false,
})

// Mock data for now
reports.value = [
	{
		id: 1,
		package: {
			id: 10,
			name: 'suspicious-pkg',
			display_name: 'Suspicious Package',
			author_slug: 'baduser',
			author_username: 'BadUser123',
		},
		reporter: {
			username: 'concerned_dev',
		},
		reason: 'Malicious code',
		details: 'This package contains code that attempts to access system files without permission.',
		status: 'pending',
		created_at: new Date(Date.now() - 3600000).toISOString(),
	},
	{
		id: 2,
		package: {
			id: 15,
			name: 'copied-lib',
			display_name: 'Copied Library',
			author_slug: 'copycat',
			author_username: 'CopyCat99',
		},
		reporter: {
			username: 'original_author',
		},
		reason: 'Copyright violation',
		details: 'This is a direct copy of my library without attribution.',
		status: 'pending',
		created_at: new Date(Date.now() - 7200000).toISOString(),
	},
	{
		id: 3,
		package: {
			id: 20,
			name: 'spam-tool',
			display_name: 'Spam Tool',
			author_slug: 'spammer',
			author_username: 'Spammer2024',
		},
		reporter: {
			username: 'moderator1',
		},
		reason: 'Spam',
		details: 'Not a real package, just advertising.',
		status: 'pending',
		created_at: new Date(Date.now() - 86400000).toISOString(),
	},
]

const filteredReports = computed(() => {
	if (filterTab.value === 'all') return reports.value
	return reports.value.filter((r) => r.status === filterTab.value)
})

const formatDate = (date) => {
	const d = new Date(date)
	const now = new Date()
	const diff = now - d

	if (diff < 3600000) return `${Math.floor(diff / 60000)} minutes ago`
	if (diff < 86400000) return `${Math.floor(diff / 3600000)} hours ago`
	return d.toLocaleDateString()
}

const getReasonColor = (reason) => {
	const colors = {
		'Malicious code': 'red',
		'Copyright violation': 'deep-orange',
		'Inappropriate content': 'orange',
		'Broken/non-functional': 'amber',
		Spam: 'grey',
		Other: 'blue-grey',
	}
	return colors[reason] || 'grey'
}

const refreshReports = async () => {
	loading.value = true
	try {
		// In production: GET /api/moderation/reports
		// const response = await api.get('/moderation/reports')
		// reports.value = response.data

		// Mock refresh
		await new Promise((resolve) => setTimeout(resolve, 1000))

		Notify.create({
			type: 'positive',
			message: 'Reports refreshed',
			position: 'top',
		})
	} catch (error) {
		console.error('Failed to refresh reports:', error)
		Notify.create({
			type: 'negative',
			message: 'Failed to refresh reports',
			position: 'top',
		})
	} finally {
		loading.value = false
	}
}

const showActionDialog = (type, report) => {
	actionDialog.value = {
		show: true,
		type,
		report,
		notes: '',
		loading: false,
	}
}

const executeAction = async () => {
	actionDialog.value.loading = true

	try {
		const { type, report, notes } = actionDialog.value

		// In production, make appropriate API calls:
		// - Remove package: DELETE /api/packages/:slug/:name
		// - Ban user: POST /api/users/:id/ban
		// - Dismiss report: PATCH /api/moderation/reports/:id

		console.log('Executing action:', { type, reportId: report.id, notes })

		// Mock success
		await new Promise((resolve) => setTimeout(resolve, 1000))

		// Update report status locally
		const idx = reports.value.findIndex((r) => r.id === report.id)
		if (idx !== -1) {
			reports.value[idx] = {
				...reports.value[idx],
				status: type === 'dismiss' ? 'dismissed' : 'resolved',
				resolved_at: new Date().toISOString(),
				resolved_by: { username: 'current_moderator' },
			}
		}

		actionDialog.value.show = false

		Notify.create({
			type: 'positive',
			message:
				type === 'remove' ? 'Package removed' : type === 'ban' ? 'User banned' : 'Report dismissed',
			position: 'top',
		})
	} catch (error) {
		console.error('Failed to execute action:', error)
		Notify.create({
			type: 'negative',
			message: 'Failed to complete action',
			position: 'top',
		})
	} finally {
		actionDialog.value.loading = false
	}
}
</script>

<style lang="scss" scoped>
.moderation-panel {
	height: 100%;
	overflow-y: auto;
}
</style>
