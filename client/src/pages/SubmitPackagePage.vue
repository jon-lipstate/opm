<template>
	<q-page class="q-pa-md">
		<div class="container" style="max-width: 800px; margin: 0 auto">
			<h1 class="text-h3 q-mb-lg">Submit Package</h1>

			<div>Submitted packages must either be Odin code, or directly support the Odin language.</div>

			<!-- Login required message -->
			<q-banner v-if="!userStore.isLoggedIn" class="bg-warning q-mb-lg">
				<template v-slot:avatar>
					<q-icon name="info" />
				</template>
				You must be logged in to submit a package.
				<template v-slot:action>
					<q-btn flat label="Login" to="/login" />
				</template>
			</q-banner>

			<q-form v-else @submit="onSubmit" @reset="onReset" class="q-gutter-md">
				<!-- Repository URL -->
				<q-input
					dense
					v-model="form.repositoryUrl"
					label="Repository URL *"
					hint="GitHub, GitLab, or other Git repository URL"
					lazy-rules
					:rules="[
						(val) => (val && val.length > 0) || 'Please enter a repository URL',
						(val) => isValidUrl(val) || 'Please enter a valid URL',
						(val) => isValidRepoUrl(val) || 'Must be a valid Git repository URL',
					]"
					@blur="onRepositoryUrlBlur"
				>
					<template v-slot:prepend>
						<q-icon name="code" />
					</template>
					<template v-slot:append>
						<q-btn
							v-if="
								form.repositoryUrl &&
								isValidUrl(form.repositoryUrl) &&
								isValidRepoUrl(form.repositoryUrl)
							"
							flat
							dense
							round
							icon="download"
							@click="fetchRepoMetadata"
							:loading="fetchingMetadata"
						>
							<q-tooltip>Fetch metadata from repository</q-tooltip>
						</q-btn>
					</template>
				</q-input>

				<!-- Package Name/Slug -->
				<q-input
					dense
					v-model="form.slug"
					label="Package Slug *"
					hint="URL-safe name (lowercase, numbers, hyphens, underscores)"
					lazy-rules
					:rules="[
						(val) => (val && val.length > 0) || 'Please enter a package slug',
						(val) =>
							/^[a-z0-9_-]+$/.test(val) ||
							'Only lowercase letters, numbers, hyphens, and underscores allowed',
						(val) => val.length >= 2 || 'Must be at least 2 characters',
						(val) => val.length <= 100 || 'Must be less than 100 characters',
					]"
					@update:model-value="(val) => (form.slug = val.toLowerCase())"
				>
					<template v-slot:prepend>
						<span class="text-caption text-grey-6">{{ userStore.slug }}/</span>
					</template>
				</q-input>

				<!-- Display Name -->
				<q-input
					dense
					v-model="form.displayName"
					label="Display Name *"
					hint="Human-readable name for your package"
					lazy-rules
					:rules="[
						(val) => (val && val.length > 0) || 'Please enter a display name',
						(val) => val.length >= 2 || 'Must be at least 2 characters',
						(val) => val.length <= 255 || 'Must be less than 255 characters',
					]"
				/>

				<!-- Description -->
				<q-input
					dense
					v-model="form.description"
					label="Description *"
					hint="Brief description of what your package does"
					type="textarea"
					rows="4"
					lazy-rules
					:rules="[
						(val) => (val && val.length > 0) || 'Please enter a description',
						(val) => val.length >= 10 || 'Must be at least 10 characters',
						(val) => val.length <= 1000 || 'Must be less than 1000 characters',
					]"
				/>

				<!-- License -->
				<q-input
					dense
					v-model="form.license"
					label="License"
					placeholder="e.g. BSD-3, MIT, Apache 2.0"
					lazy-rules
					:rules="[
						(val) => !val || val.length <= 100 || 'License must be less than 100 characters',
					]"
				>
					<template v-slot:prepend>
						<q-icon name="gavel" />
					</template>
				</q-input>

				<!-- Type -->
				<div>
					<div class="text-subtitle2 q-mb-sm">Package Type *</div>
					<q-option-group dense v-model="form.type" :options="typeOptions" inline />
				</div>

				<!-- Status -->
				<div>
					<div class="text-subtitle2 q-mb-sm">Development Status *</div>
					<q-option-group dense v-model="form.status" :options="statusOptions" inline />
				</div>

				<!-- Tags -->
				<div>
					<div class="text-subtitle2 q-mb-sm">Tags</div>
					<div class="text-caption text-grey-6 q-mb-sm">
						Add relevant tags to help others discover your package
					</div>

					<div class="row q-gutter-sm q-mb-sm">
						<q-chip
							v-for="tag in selectedTags"
							:key="tag.id"
							removable
							@remove="removeTag(tag)"
							color="primary"
							text-color="white"
						>
							{{ tag.name }}
						</q-chip>
					</div>

					<q-select
						dense
						v-model="tagInput"
						:options="filteredTags"
						option-label="name"
						option-value="id"
						label="Add tags"
						use-input
						input-debounce="300"
						@filter="filterTags"
						@update:model-value="addTag"
						clearable
						hide-dropdown-icon
					>
						<template v-slot:option="scope">
							<q-item v-bind="scope.itemProps">
								<q-item-section>
									<q-item-label>{{ scope.opt.name }}</q-item-label>
									<q-item-label caption>Used {{ scope.opt.usage_count }} times</q-item-label>
								</q-item-section>
							</q-item>
						</template>

						<template v-slot:no-option="scope">
							<q-item>
								<q-item-section>
									<q-item-label>
										<span v-if="scope.inputValue"
											>Press Enter to create "{{ scope.inputValue }}"</span
										>
										<span v-else>Start typing to search or create tags</span>
									</q-item-label>
								</q-item-section>
							</q-item>
						</template>
					</q-select>
				</div>

				<!-- Submit/Reset buttons -->
				<div class="q-mt-xl">
					<q-btn
						label="Submit Package"
						type="submit"
						color="primary"
						:loading="submitting"
						:disable="!isFormValid"
					/>
					<q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
				</div>
			</q-form>
		</div>
	</q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from 'stores/user-store'
import { useApiStore } from 'stores/api-store'
import { Notify } from 'quasar'

const router = useRouter()
const userStore = useUserStore()
const apiStore = useApiStore()

// Form data
const form = ref({
	slug: '',
	displayName: '',
	description: '',
	repositoryUrl: '',
	license: '',
	type: 'project',
	status: 'in_work',
})

const submitting = ref(false)
const fetchingMetadata = ref(false)
const allTags = ref([])
const filteredTags = ref([])
const selectedTags = ref([])
const tagInput = ref(null)

const typeOptions = [
	{ label: 'Library', value: 'library', color: 'primary' },
	{ label: 'Project', value: 'project', color: 'secondary' },
]

const statusOptions = [
	{ label: 'Work in Progress', value: 'in_work', color: 'secondary' },
	{ label: 'Done', value: 'ready', color: 'positive' },
]

// Computed
const isFormValid = computed(() => {
	return (
		form.value.slug &&
		form.value.displayName &&
		form.value.description &&
		form.value.repositoryUrl &&
		isValidUrl(form.value.repositoryUrl) &&
		isValidRepoUrl(form.value.repositoryUrl) &&
		/^[a-z0-9_-]+$/.test(form.value.slug)
	)
})

// Methods
const isValidUrl = (url) => {
	try {
		new URL(url)
		return true
	} catch {
		return false
	}
}

const isValidRepoUrl = (url) => {
	const repoPatterns = [
		/^https?:\/\/(www\.)?github\.com\/.+\/.+/,
		/^https?:\/\/(www\.)?gitlab\.com\/.+\/.+/,
		/^https?:\/\/(www\.)?bitbucket\.org\/.+\/.+/,
		/^https?:\/\/.+\/.+\/.+/, // Generic pattern for self-hosted
	]

	return repoPatterns.some((pattern) => pattern.test(url))
}

const loadTags = async () => {
	try {
		const tags = await apiStore.fetchTags({ limit: 100 })
		allTags.value = tags
		filteredTags.value = allTags.value
	} catch (_error) {
		// Error already handled by apiStore
	}
}

const filterTags = (val, update) => {
	update(() => {
		const needle = val.toLowerCase()
		filteredTags.value = allTags.value
			.filter((tag) => !selectedTags.value.some((t) => t.id === tag.id))
			.filter((tag) => tag.name.toLowerCase().includes(needle))

		// Add option to create new tag if it doesn't exist
		if (val && !filteredTags.value.some((tag) => tag.name.toLowerCase() === needle)) {
			filteredTags.value.unshift({
				id: 'new',
				name: val,
				usage_count: 0,
				isNew: true,
			})
		}
	})
}

const addTag = (tag) => {
	if (!tag) return

	if (tag.isNew) {
		// Create a temporary tag that will be created on the server
		const newTag = {
			id: `new-${Date.now()}`,
			name: tag.name,
			isNew: true,
		}
		selectedTags.value.push(newTag)
	} else {
		selectedTags.value.push(tag)
	}

	tagInput.value = null
}

const removeTag = (tag) => {
	const index = selectedTags.value.indexOf(tag)
	if (index > -1) {
		selectedTags.value.splice(index, 1)
	}
}

const onSubmit = async () => {
	if (!isFormValid.value) return

	submitting.value = true

	try {
		// Prepare tag IDs (filter out new tags)
		const tagIds = selectedTags.value
			.filter((tag) => !tag.isNew && typeof tag.id === 'number')
			.map((tag) => tag.id)

		const response = await apiStore.createPackage({
			slug: form.value.slug,
			display_name: form.value.displayName,
			description: form.value.description,
			repository_url: form.value.repositoryUrl,
			license: form.value.license || undefined,
			type: form.value.type,
			status: form.value.status,
			tag_ids: tagIds,
		})

		// Add new tags after package creation
		const newTags = selectedTags.value.filter((tag) => tag.isNew)
		for (const tag of newTags) {
			try {
				await apiStore.addPackageTag(response.id, tag.name)
			} catch (error) {
				console.error('Failed to add tag:', tag.name, error)
			}
		}

		// Success notification is already shown by apiStore.createPackage

		// Redirect to the new package page
		router.push(`/packages/${userStore.slug}/${response.slug}`)
	} catch (error) {
		console.error('Failed to submit package:', error)

		let errorMessage = 'Failed to submit package'
		if (error.response?.status === 409) {
			errorMessage = 'A package with this slug already exists'
		} else if (error.response?.data?.message) {
			errorMessage = error.response.data.message
		}

		Notify.create({
			type: 'negative',
			message: errorMessage,
			position: 'top',
		})
	} finally {
		submitting.value = false
	}
}

const onReset = () => {
	form.value = {
		slug: '',
		displayName: '',
		description: '',
		repositoryUrl: '',
		license: '',
		type: 'library',
		status: 'in_work',
	}
	selectedTags.value = []
	tagInput.value = null
}

const fetchRepoMetadata = async () => {
	if (!form.value.repositoryUrl || !isValidUrl(form.value.repositoryUrl)) {
		return
	}

	fetchingMetadata.value = true
	try {
		const metadata = await apiStore.fetchRepositoryMetadata(form.value.repositoryUrl)
		if (metadata) {
			// Always update fields when fetching metadata
			if (metadata.slug) {
				form.value.slug = metadata.slug
			}
			if (metadata.display_name) {
				form.value.displayName = metadata.display_name
			}
			if (metadata.description) {
				form.value.description = metadata.description
			}
			if (metadata.license) {
				form.value.license = metadata.license
			}

			Notify.create({
				type: 'positive',
				message: 'Repository metadata fetched successfully',
				position: 'top',
			})
		}
	} catch (error) {
		console.error('Failed to fetch repository metadata:', error)
		// Error is already handled by apiStore
	} finally {
		fetchingMetadata.value = false
	}
}

const onRepositoryUrlBlur = () => {
	// Auto-fetch metadata when repository URL is entered if other fields are empty
	if (
		form.value.repositoryUrl &&
		isValidUrl(form.value.repositoryUrl) &&
		isValidRepoUrl(form.value.repositoryUrl) &&
		!form.value.displayName &&
		!form.value.description
	) {
		fetchRepoMetadata()
	}
}

// Lifecycle
onMounted(() => {
	if (!userStore.isLoggedIn) {
		return
	}

	loadTags()
})
</script>

<style lang="scss" scoped>
.container {
	:deep(.q-field__label) {
		color: var(--odin-blue);
	}

	:deep(.q-field--focused .q-field__label) {
		color: var(--odin-blue-dark);
	}
}
</style>
