<template>
	<q-page>
		<!-- Hero Section -->
		<div class="odin-hero" :class="hasSearched || packages.length > 0 ? 'q-py-md' : 'q-py-xl'">
			<div class="container q-px-md" style="max-width: 1200px; margin: 0 auto">
				<div class="text-center">
					<h1
						class="text-h2 text-weight-bold q-my-none"
						v-show="!hasSearched && packages.length === 0"
					>
						Odin Package Registry
					</h1>
					<p
						class="text-h6 text-odin-light q-mt-md"
						style="opacity: 0.9"
						v-show="!hasSearched && packages.length === 0"
					>
						Discover and share Odin libraries and projects
					</p>

					<!-- Search Box -->
					<div
						:class="hasSearched || packages.length > 0 ? '' : 'q-mt-xl'"
						style="max-width: 600px; margin: 0 auto"
					>
						<q-input
							ref="searchInput"
							v-model="searchQuery"
							outlined
							placeholder="Search packages or topics"
							bg-color="white"
							class="search-input"
							@keyup.enter="search"
						>
							<template v-slot:append>
								<span class="text-caption text-grey-5 q-mr-sm">Ctrl+/</span>
								<q-btn icon="search" flat dense @click="search" />
							</template>
						</q-input>
					</div>
				</div>
			</div>
		</div>

		<!-- Filters -->
		<div class="bg-grey-2 q-py-md">
			<div class="container q-px-md" style="max-width: 1200px; margin: 0 auto">
				<div class="row q-gutter-md items-center">
					<div>
						<q-btn-toggle
							v-model="filters.type"
							toggle-color="primary"
							:options="[
								{ label: 'All', value: 'all' },
								{ label: 'Libraries', value: 'library' },
								{ label: 'Projects', value: 'project' },
							]"
							@update:model-value="
								() => {
									hasSearched = true
									loadPackages()
								}
							"
						/>
					</div>

					<div>
						<q-btn-toggle
							v-model="filters.status"
							toggle-color="primary"
							:options="[
								{ label: 'All', value: 'all' },
								{ label: 'Ready', value: 'ready' },
								{ label: 'In Work', value: 'in_work' },
							]"
							@update:model-value="
								() => {
									hasSearched = true
									loadPackages()
								}
							"
						/>
					</div>

					<q-space />
				</div>
			</div>
		</div>

		<!-- Package List -->
		<div class="container q-pa-md" style="max-width: 1200px; margin: 0 auto">
			<div v-if="loading" class="text-center q-py-xl">
				<q-spinner size="50px" color="primary" />
			</div>

			<div v-else-if="packages.length === 0" class="text-center q-py-xl">
				<q-icon name="inventory_2" size="64px" color="grey-5" />
				<p class="text-h6 text-grey-7 q-mt-md">No packages found</p>
				<p class="text-body2 text-grey-6">Try adjusting your filters or search terms</p>
			</div>

			<div v-else class="row q-col-gutter-md">
				<div v-for="pkg in packages" :key="pkg.id" class="col-12 col-md-6 col-lg-4">
					<q-card class="full-height odin-card relative-position">
						<!-- Flag indicator in top right -->
						<q-btn
							v-if="pkg.active_reports_count > 0"
							flat
							round
							dense
							size="md"
							icon="warning"
							color="warning"
							class="absolute-top-right q-ma-sm"
							style="z-index: 1; cursor: default"
						>
							<q-badge color="negative" floating rounded :label="pkg.active_reports_count" />
							<q-tooltip class="text-subtitle1"
								>{{ pkg.active_reports_count }} pending moderation report{{
									pkg.active_reports_count > 1 ? 's' : ''
								}}</q-tooltip
							>
						</q-btn>

						<q-card-section class="q-pb-xs">
							<router-link
								:to="`/packages/${pkg.author?.slug || pkg.author?.username}/${pkg.slug}`"
								class="text-h6 text-no-underline"
								style="color: var(--odin-blue)"
							>
								{{ pkg.display_name }}
							</router-link>

							<!-- Type and Status badges right under title -->
							<div class="row items-center q-mt-xs q-gutter-xs">
								<q-badge
									:color="pkg.type === 'library' ? 'primary' : 'orange'"
									:outline="pkg.type === 'project'"
									class="text-caption"
								>
									{{ pkg.type }}
								</q-badge>
								<q-badge
									:color="pkg.status === 'ready' ? 'positive' : 'amber'"
									:outline="pkg.status === 'in_work'"
									class="text-caption"
								>
									{{ pkg.status === 'ready' ? 'ready' : 'wip' }}
								</q-badge>
							</div>
						</q-card-section>

						<q-card-section class="q-pt-none">
							<p class="text-body2 text-grey-8 q-mb-sm">
								{{
									pkg.description && pkg.description.length > 128
										? pkg.description.substring(0, 128) + '...'
										: pkg.description
								}}
							</p>

							<div class="row items-center text-caption text-grey-6">
								<q-avatar size="20px" class="q-mr-xs">
									<img v-if="pkg.author?.avatar_url" :src="pkg.author.avatar_url" />
									<q-icon v-else name="account_circle" />
								</q-avatar>
								<span>{{ pkg.author?.display_name || pkg.author?.username || 'Unknown' }}</span>
							</div>

							<div v-if="pkg.tags?.length" class="q-mt-sm">
								<q-chip
									v-for="tag in pkg.tags"
									:key="tag.id"
									size="sm"
									outline
									color="grey-7"
									dense
									class="q-mr-xs"
								>
									{{ tag.name }}
								</q-chip>
							</div>
						</q-card-section>
					</q-card>
				</div>
			</div>

			<!-- Load More -->
			<div v-if="hasMore" class="text-center q-mt-lg">
				<q-btn label="Load More" color="primary" outline @click="loadMore" :loading="loading" />
			</div>
		</div>
	</q-page>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useApiStore } from 'stores/api-store'

const apiStore = useApiStore()

const searchInput = ref(null)
const searchQuery = ref('')
const packages = ref([])
const loading = ref(false)
const hasMore = ref(false)
const hasSearched = ref(false)

const filters = ref({
	type: 'all',
	status: 'all',
	limit: 12,
	offset: 0,
})

const search = () => {
	hasSearched.value = true
	filters.value.offset = 0
	packages.value = []
	loadPackages()
}

const loadPackages = async () => {
	loading.value = true

	try {
		let data

		if (searchQuery.value) {
			// Use search endpoint when searching
			data = await apiStore.searchPackages(searchQuery.value, {
				limit: filters.value.limit,
				offset: filters.value.offset,
			})
		} else {
			// Use regular listing endpoint with filters
			const params = {
				limit: filters.value.limit,
				offset: filters.value.offset,
			}

			if (filters.value.type !== 'all') {
				params.type = filters.value.type
			}
			if (filters.value.status !== 'all') {
				params.status = filters.value.status
			}

			data = await apiStore.fetchPackages(params)
		}

		if (filters.value.offset === 0) {
			packages.value = data
		} else {
			packages.value.push(...data)
		}

		hasMore.value = data.length === filters.value.limit
	} catch (_error) {
		// Error already handled by apiStore
	} finally {
		loading.value = false
	}
}

const loadMore = () => {
	filters.value.offset += filters.value.limit
	loadPackages()
}

// Keyboard shortcut handler
const handleKeydown = (e) => {
	// Ctrl+/ or Cmd+/ to focus search
	if ((e.ctrlKey || e.metaKey) && e.key === '/') {
		e.preventDefault()
		searchInput.value?.focus()
	}
}

onMounted(() => {
	// Load packages from API
	loadPackages()
	window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
	window.removeEventListener('keydown', handleKeydown)
})
</script>
