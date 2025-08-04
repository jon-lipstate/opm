<template>
	<q-page class="q-pa-md">
		<div class="container">
			<div class="row items-center q-mb-lg">
				<h4 class="q-ma-none">My Packages</h4>
				<q-space />
				<q-btn color="primary" label="Submit Package" icon="add" to="/submit" />
			</div>

			<q-card>
				<q-table
					:rows="packages"
					:columns="columns"
					row-key="id"
					:loading="loading"
					:pagination="{ rowsPerPage: 20 }"
					flat
					:no-data-label="loading ? 'Loading packages...' : 'No packages found'"
				>
					<!-- Package Name/Title -->
					<template v-slot:body-cell-name="props">
						<q-td :props="props">
							<router-link
								:to="`/packages/${userStore.user.slug}/${props.row.slug}`"
								class="text-primary text-weight-medium"
								style="text-decoration: none"
							>
								{{ props.row.display_name }}
							</router-link>
							<div class="text-caption text-grey-6">{{ props.row.slug }}</div>
						</q-td>
					</template>

					<!-- Type -->
					<template v-slot:body-cell-type="props">
						<q-td :props="props">
							<q-chip
								:color="props.row.type === 'library' ? 'primary' : 'orange'"
								:outline="props.row.type === 'project'"
								text-color="white"
								size="sm"
								dense
							>
								{{ props.row.type }}
							</q-chip>
						</q-td>
					</template>

					<!-- Status -->
					<template v-slot:body-cell-status="props">
						<q-td :props="props">
							<q-chip :color="getStatusColor(props.row.status)" text-color="white" size="sm" dense>
								{{ formatStatus(props.row.status) }}
							</q-chip>
						</q-td>
					</template>

					<!-- Actions -->
					<template v-slot:body-cell-actions="props">
						<q-td :props="props" class="q-gutter-xs">
							<q-btn
								flat
								dense
								round
								icon="edit"
								size="sm"
								color="primary"
								@click="editPackage(props.row)"
							>
								<q-tooltip>Edit Package</q-tooltip>
							</q-btn>
							<q-btn
								flat
								dense
								round
								icon="delete"
								size="sm"
								color="negative"
								@click="confirmDelete(props.row)"
							>
								<q-tooltip>Delete Package</q-tooltip>
							</q-btn>
						</q-td>
					</template>
				</q-table>
			</q-card>
		</div>

		<!-- Edit Dialog -->
		<q-dialog v-model="showEditDialog">
			<q-card style="min-width: 500px">
				<q-card-section>
					<div class="text-h6">Edit Package</div>
				</q-card-section>

				<q-card-section>
					<q-form @submit="updatePackage">
						<q-input
							v-model="editForm.repository_url"
							label="Repository URL"
							outlined
							:rules="[(val) => !!val || 'Repository URL is required']"
						/>

						<q-input
							v-model="editForm.display_name"
							label="Display Name"
							outlined
							class="q-mb-md"
							:rules="[(val) => !!val || 'Display name is required']"
						/>

						<q-input
							v-model="editForm.description"
							label="Description"
							outlined
							type="textarea"
							rows="3"
							class="q-mb-md"
							:rules="[(val) => !!val || 'Description is required']"
						/>

						<q-select
							v-model="editForm.type"
							label="Type"
							outlined
							:options="typeOptions"
							emit-value
							map-options
							class="q-mb-md"
						/>

						<q-select
							v-model="editForm.status"
							label="Status"
							outlined
							:options="statusOptions"
							emit-value
							map-options
							class="q-mb-md"
						/>

						<q-input
							v-model="editForm.license"
							label="License"
							outlined
							placeholder="e.g. BSD-3, MIT, Apache 2.0"
						/>
					</q-form>
				</q-card-section>

				<q-card-actions align="right">
					<q-btn flat label="Cancel" v-close-popup />
					<q-btn color="primary" label="Update" @click="updatePackage" :loading="editLoading" />
				</q-card-actions>
			</q-card>
		</q-dialog>

		<!-- Delete Confirmation Dialog -->
		<q-dialog v-model="showDeleteDialog">
			<q-card>
				<q-card-section>
					<div class="text-h6">Delete Package</div>
				</q-card-section>

				<q-card-section>
					Are you sure you want to delete "{{ packageToDelete?.display_name }}"? This action cannot
					be undone.
				</q-card-section>

				<q-card-actions align="right">
					<q-btn flat label="Cancel" v-close-popup />
					<q-btn color="negative" label="Delete" @click="deletePackage" :loading="deleteLoading" />
				</q-card-actions>
			</q-card>
		</q-dialog>
	</q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useUserStore } from 'src/stores/user-store'
import { useApiStore } from 'src/stores/api-store'

const $q = useQuasar()
const router = useRouter()
const userStore = useUserStore()
const apiStore = useApiStore()

// Data
const packages = ref([])
const loading = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const editLoading = ref(false)
const deleteLoading = ref(false)
const packageToDelete = ref(null)

const editForm = ref({
	id: null,
	display_name: '',
	description: '',
	type: '',
	status: '',
	repository_url: '',
	license: '',
})

// Table columns
const columns = [
	{
		name: 'name',
		label: 'Package',
		field: 'display_name',
		align: 'left',
		sortable: true,
	},
	{
		name: 'type',
		label: 'Type',
		field: 'type',
		align: 'center',
		sortable: true,
	},
	{
		name: 'status',
		label: 'Status',
		field: 'status',
		align: 'center',
		sortable: true,
	},
	{
		name: 'views',
		label: 'Views',
		field: 'view_count',
		align: 'center',
		sortable: true,
		format: (val) => val.toLocaleString(),
	},
	{
		name: 'bookmarks',
		label: 'Bookmarks',
		field: 'bookmark_count',
		align: 'center',
		sortable: true,
		format: (val) => val.toLocaleString(),
	},
	{
		name: 'actions',
		label: 'Actions',
		align: 'center',
	},
]

const typeOptions = [
	{ label: 'Library', value: 'library' },
	{ label: 'Project', value: 'project' },
]

const statusOptions = [
	{ label: 'In Work', value: 'in_work' },
	{ label: 'Ready', value: 'ready' },
	{ label: 'Archived', value: 'archived' },
	{ label: 'Abandoned', value: 'abandoned' },
]

// Methods
const loadPackages = async () => {
	loading.value = true
	try {
		const result = await apiStore.fetchUserPackages()
		packages.value = result || []
	} catch (error) {
		console.error('Failed to load packages:', error)
		// Only redirect to login if we get a 401 Unauthorized
		if (error.response?.status === 401) {
			router.push('/login')
		} else {
			$q.notify({
				type: 'negative',
				message: 'Failed to load packages',
			})
		}
	} finally {
		loading.value = false
	}
}

const getStatusColor = (status) => {
	const colors = {
		in_work: 'warning',
		ready: 'positive',
		archived: 'grey',
		abandoned: 'negative',
	}
	return colors[status] || 'grey'
}

const formatStatus = (status) => {
	return status
		.split('_')
		.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
		.join(' ')
}

const editPackage = (pkg) => {
	editForm.value = {
		id: pkg.id,
		display_name: pkg.display_name,
		description: pkg.description,
		type: pkg.type,
		status: pkg.status,
		repository_url: pkg.repository_url,
		license: pkg.license || '',
	}
	showEditDialog.value = true
}

const updatePackage = async () => {
	editLoading.value = true
	try {
		await apiStore.updatePackage(editForm.value.id, {
			display_name: editForm.value.display_name,
			description: editForm.value.description,
			type: editForm.value.type,
			status: editForm.value.status,
			repository_url: editForm.value.repository_url,
			license: editForm.value.license || undefined,
		})

		$q.notify({
			type: 'positive',
			message: 'Package updated successfully',
		})

		showEditDialog.value = false
		loadPackages() // Reload the list
	} catch (error) {
		console.error('Failed to update package:', error)
		$q.notify({
			type: 'negative',
			message: 'Failed to update package',
		})
	} finally {
		editLoading.value = false
	}
}

const confirmDelete = (pkg) => {
	packageToDelete.value = pkg
	showDeleteDialog.value = true
}

const deletePackage = async () => {
	if (!packageToDelete.value) return

	deleteLoading.value = true
	try {
		await apiStore.deletePackage(packageToDelete.value.id)

		$q.notify({
			type: 'positive',
			message: 'Package deleted successfully',
		})

		showDeleteDialog.value = false
		loadPackages() // Reload the list
	} catch (error) {
		console.error('Failed to delete package:', error)
		$q.notify({
			type: 'negative',
			message: 'Failed to delete package',
		})
	} finally {
		deleteLoading.value = false
		packageToDelete.value = null
	}
}

onMounted(async () => {
	// Wait for user store to initialize if it's still loading
	if (userStore.isLoading) {
		// Wait for the user store to finish loading
		const unwatch = watch(
			() => userStore.isLoading,
			(newVal) => {
				if (!newVal) {
					unwatch()
					loadPackages()
				}
			},
		)
	} else {
		// User store already loaded
		loadPackages()
	}
})
</script>

<style scoped>
.container {
	max-width: 1200px;
	margin: 0 auto;
}
</style>
