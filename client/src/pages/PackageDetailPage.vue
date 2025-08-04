<template>
	<q-page>
		<div class="container q-pa-md" style="max-width: 1200px; margin: 0 auto">
			<div v-if="loading" class="text-center q-py-xl">
				<q-spinner size="50px" color="primary" />
			</div>

			<div v-else-if="!pkg" class="text-center q-py-xl">
				<q-icon name="error_outline" size="64px" color="negative" />
				<p class="text-h6 q-mt-md">Package not found</p>
				<q-btn label="Go Home" color="primary" to="/" class="q-mt-md" />
			</div>

			<div v-else>
				<!-- Package Header -->
				<div class="row items-start q-col-gutter-lg">
					<!-- Main content - order-last only on small screens -->
					<div class="col-12 col-md-8" :class="$q.screen.lt.md ? 'order-last' : ''">
						<div class="row items-center q-gutter-sm q-mb-sm">
							<h1 class="text-h3 q-my-none">{{ pkg.display_name }}</h1>
							<q-badge
								:color="pkg.type === 'library' ? 'primary' : 'orange'"
								:outline="pkg.type === 'project'"
								class="text-body2"
							>
								{{ pkg.type }}
							</q-badge>
							<q-badge
								:color="pkg.status === 'ready' ? 'positive' : 'amber'"
								:outline="pkg.status === 'in_work'"
								class="text-body2"
							>
								{{ pkg.status === 'ready' ? 'ready' : 'wip' }}
							</q-badge>
						</div>

						<p class="text-h6 text-grey-7 q-mb-md">{{ pkg.description }}</p>

						<!-- Author info -->
						<div class="row items-center text-body2 text-grey-6 q-mb-lg">
							<div class="row items-center q-gutter-sm">
								<q-avatar size="24px">
									<img v-if="pkg.author?.avatar_url" :src="pkg.author.avatar_url" />
									<q-icon v-else name="account_circle" />
								</q-avatar>
								<span
									>by
									<router-link :to="`/users/${pkg.author?.slug}`" class="text-primary">
										{{ pkg.author?.display_name || pkg.author?.username }}
									</router-link>
									<q-icon
										v-if="pkg.author?.reputation_rank"
										:name="getReputationIcon(pkg.author.reputation_rank)"
										:color="getReputationColor(pkg.author.reputation_rank)"
										size="16px"
										class="q-ml-xs"
									>
										<q-tooltip>
											{{ pkg.author.reputation_rank }} ({{ pkg.author.reputation }} points)
										</q-tooltip>
									</q-icon>
								</span>
								<span>•</span>
								<q-icon name="bookmark" size="16px" />
								<span>{{ pkg.bookmark_count || 0 }}</span>
								<span>•</span>
								<q-icon name="visibility" size="16px" />
								<span>{{ pkg.view_count || 0 }}</span>
								<span>•</span>
								<span>Updated {{ timeAgo(pkg.updated_at) }}</span>
							</div>

							<q-space />

							<!-- Flag button (subtle) -->
							<q-btn
								v-if="userStore.isLoggedIn"
								flat
								dense
								round
								size="sm"
								icon="flag"
								color="grey-6"
								@click="showFlagDialog = true"
							>
								Report Issue
								<q-tooltip class="text-subtitle2"
									>Report package as inappropriate, malware or other security concerns</q-tooltip
								>
							</q-btn>
						</div>

						<!-- Tags -->
						<div class="q-mb-lg">
							<div class="row items-center q-mb-sm">
								<span class="text-subtitle2 text-grey-7">Tags</span>
								<q-space />
								<q-btn
									v-if="userStore.isLoggedIn"
									flat
									dense
									size="sm"
									icon="add"
									label="Add Tag"
									color="primary"
									@click="showAddTagDialog = true"
								/>
							</div>

							<div v-if="!pkg.tags?.length" class="text-caption text-grey-6">
								No tags yet. Be the first to add one!
							</div>

							<div v-else class="row q-gutter-xs">
								<div v-for="tag in sortedTags" :key="tag.id" class="tag-with-voting">
									<q-chip
										:outline="tag.net_score < 5"
										:color="getTagColor(tag)"
										clickable
										@click="$router.push({ path: '/', query: { tag: tag.name } })"
										class="q-ma-none"
									>
										{{ tag.name }}
										<q-badge
											v-if="tag.net_score !== 0"
											:color="tag.net_score > 0 ? 'positive' : 'negative'"
											floating
											transparent
											:label="tag.net_score > 0 ? `+${tag.net_score}` : tag.net_score"
										/>
									</q-chip>

									<!-- Voting buttons -->
									<div v-if="userStore.isLoggedIn" class="tag-voting-buttons">
										<q-btn
											flat
											dense
											round
											size="xs"
											icon="arrow_upward"
											:color="tag.user_vote > 0 ? 'positive' : 'grey'"
											@click="voteOnTag(tag, 1)"
											:loading="tag.voting"
										>
											<q-tooltip>Upvote this tag</q-tooltip>
										</q-btn>
										<q-btn
											flat
											dense
											round
											size="xs"
											icon="arrow_downward"
											:color="tag.user_vote < 0 ? 'negative' : 'grey'"
											@click="voteOnTag(tag, -1)"
											:loading="tag.voting"
										>
											<q-tooltip>Downvote this tag</q-tooltip>
										</q-btn>
									</div>
								</div>
							</div>
						</div>

						<!-- README Content -->
						<PackageReadme :content="readmeContent" :loading="readmeLoading" :error="readmeError" />
					</div>

					<!-- Sidebar - order-first only on small screens -->
					<div class="col-12 col-md-4" :class="$q.screen.lt.md ? 'order-first' : ''">
						<!-- Actions Card -->
						<q-card class="q-mb-md">
							<q-card-section>
								<div class="q-gutter-sm">
									<!-- Repository Link -->
									<q-btn
										:href="pkg.repository_url"
										target="_blank"
										color="primary"
										class="full-width"
										icon="code"
										label="View Repository"
										no-caps
									/>

									<!-- Bookmark Button -->
									<q-btn
										:color="isBookmarked ? 'secondary' : 'grey'"
										class="full-width"
										:icon="isBookmarked ? 'bookmark' : 'bookmark_border'"
										:label="isBookmarked ? 'Bookmarked' : 'Bookmark'"
										no-caps
										@click="toggleBookmark"
										:loading="bookmarkLoading"
									/>

									<!-- Clone Command -->
									<div class="q-mt-md">
										<div class="text-caption text-grey-6 q-mb-xs">Clone command:</div>
										<q-input
											:model-value="`git clone ${pkg.repository_url}`"
											readonly
											dense
											outlined
											class="clone-command"
										>
											<template v-slot:append>
												<q-btn flat round dense icon="content_copy" @click="copyCloneCommand">
													<q-tooltip>Copy command</q-tooltip>
												</q-btn>
											</template>
										</q-input>
									</div>
								</div>
							</q-card-section>
						</q-card>

						<!-- Package Info Card -->
						<q-card class="q-mb-md">
							<q-card-section>
								<div class="text-h6 q-mb-md">Package Info</div>
								<q-list dense>
									<q-item>
										<q-item-section avatar>
											<q-icon name="event" />
										</q-item-section>
										<q-item-section>
											<q-item-label>Created</q-item-label>
										</q-item-section>
										<q-item-section side>
											<q-item-label caption>{{ formatDate(pkg.created_at) }}</q-item-label>
										</q-item-section>
									</q-item>

									<q-item>
										<q-item-section avatar>
											<q-icon name="code" />
										</q-item-section>
										<q-item-section>
											<q-item-label>Repository</q-item-label>
										</q-item-section>
										<q-item-section side>
											<q-item-label caption>{{ getRepoType(pkg.repository_url) }}</q-item-label>
										</q-item-section>
									</q-item>

									<q-item v-if="pkg.license">
										<q-item-section avatar>
											<q-icon name="gavel" />
										</q-item-section>
										<q-item-section>
											<q-item-label>License</q-item-label>
										</q-item-section>
										<q-item-section side>
											<q-item-label caption>{{
												pkg.license && pkg.license.length > 16
													? pkg.license.substring(0, 16) + '...'
													: pkg.license
											}}</q-item-label>
										</q-item-section>
									</q-item>
								</q-list>
							</q-card-section>
						</q-card>

						<!-- Active Reports Card (publicly visible) -->
						<q-card
							v-if="activeReports.length > 0"
							class="q-mb-md"
							style="border: 2px solid var(--q-warning)"
						>
							<q-card-section>
								<div class="row items-center q-mb-md">
									<q-icon name="warning" size="24px" class="q-mr-sm" />
									<span class="text-h6">Moderation Reports</span>
								</div>

								<div class="text-body2 q-mb-md">
									{{ activeReports.length }} pending report{{ activeReports.length > 1 ? 's' : '' }}
								</div>

								<q-list separator>
									<q-item v-for="report in activeReports" :key="report.id" class="q-pa-none">
										<q-item-section>
											<q-item-label>
												<div class="row items-center justify-between">
													<div class="row items-center">
														<q-chip
															:color="getReportSeverityColor(report.reason)"
															text-color="white"
															size="sm"
															dense
														>
															{{ report.reason }}
														</q-chip>
														<q-btn
															v-if="userStore.user?.id === report.user_id"
															flat
															dense
															label="Retract"
															size="sm"
															color="negative"
															class="q-ml-md"
															@click="retractReport(report.id)"
															:loading="report.retracting"
														/>
													</div>
													<span class="text-caption text-grey-6">
														{{ timeAgo(report.created_at) }}
													</span>
												</div>
											</q-item-label>
											<q-item-label
												v-if="report.details && userStore.isLoggedIn"
												caption
												class="q-mt-xs"
											>
												"{{ report.details }}"
											</q-item-label>
										</q-item-section>
									</q-item>
								</q-list>

								<div class="text-caption q-mt-md">
									<q-icon name="info" size="14px" />
									Reports shown have not been verified. Valid reports result in package removal
								</div>
							</q-card-section>
						</q-card>
					</div>
				</div>
			</div>
		</div>

		<!-- Flag Dialog -->
		<q-dialog v-model="showFlagDialog">
			<q-card style="min-width: 400px">
				<q-card-section>
					<div class="text-h6">Flag Package</div>
				</q-card-section>

				<q-card-section>
					<q-select
						v-model="flagReason"
						:options="flagReasons"
						label="Reason"
						outlined
						class="q-mb-md"
					/>
					<q-input
						v-model="flagDetails"
						type="textarea"
						label="Additional details (optional)"
						outlined
						rows="4"
					/>
				</q-card-section>

				<q-card-actions align="right">
					<q-btn flat label="Cancel" v-close-popup />
					<q-btn
						color="negative"
						label="Submit Report"
						@click="submitFlag"
						:loading="flagLoading"
					/>
				</q-card-actions>
			</q-card>
		</q-dialog>

		<!-- Add Tag Dialog -->
		<q-dialog v-model="showAddTagDialog">
			<q-card style="min-width: 400px">
				<q-card-section>
					<div class="text-h6">Add Tag</div>
				</q-card-section>

				<q-card-section>
					<q-input v-model="newTagName" label="Tag name" outlined placeholder="Enter a tag name" />

					<div class="text-caption text-grey-6 q-mt-sm">You can vote ±1 on tags</div>
				</q-card-section>

				<q-card-actions align="right">
					<q-btn flat label="Cancel" v-close-popup />
					<q-btn
						color="primary"
						label="Add Tag"
						@click="addTag"
						:loading="addingTag"
						:disable="!newTagName"
					/>
				</q-card-actions>
			</q-card>
		</q-dialog>
	</q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from 'stores/user-store'
import { useApiStore } from 'stores/api-store'
import { Notify, copyToClipboard } from 'quasar'
import PackageReadme from 'src/components/PackageReadme.vue'
import { timeAgo } from 'src/utils/utils.js'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const apiStore = useApiStore()

const loading = ref(false)
const pkg = ref(null)
const readmeContent = ref('')
const readmeLoading = ref(false)
const readmeError = ref(false)
const bookmarkLoading = ref(false)
const showFlagDialog = ref(false)
const flagLoading = ref(false)
const flagReason = ref(null)
const flagDetails = ref('')
const activeReports = ref([])
const showAddTagDialog = ref(false)
const newTagName = ref('')
const availableTags = ref([])
const addingTag = ref(false)

const flagReasons = [
	'Malicious code',
	'Copyright violation',
	'Inappropriate content',
	'Broken/non-functional',
	'Spam',
	'Other',
]

const isBookmarked = computed(() => pkg.value?.is_bookmarked || false)

const sortedTags = computed(() => {
	if (!pkg.value?.tags) return []
	return [...pkg.value.tags].sort((a, b) => (b.net_score || 0) - (a.net_score || 0))
})

const formatDate = (date) => {
	return new Date(date).toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
	})
}

const loadPackage = async () => {
	loading.value = true
	try {
		const packageData = await apiStore.fetchPackage(route.params.userSlug, route.params.packageSlug)

		// Add voting state to tags
		if (packageData.tags) {
			packageData.tags = packageData.tags.map((tag) => ({
				...tag,
				voting: false,
			}))
		}

		pkg.value = packageData
		// allow these to come in async without disrupting ui:
		loadReadme()
		loadFlags()
	} catch (error) {
		console.error('Failed to load package:', error)
		pkg.value = null
	} finally {
		loading.value = false
	}
}

const loadFlags = async () => {
	// Load active flags
	try {
		const flags = await apiStore.fetchPackageFlags(pkg.value.id)
		// Add retracting state to each flag
		activeReports.value = flags.map((flag) => ({
			...flag,
			retracting: false,
		}))
	} catch (error) {
		console.error('Failed to load flags:', error)
		activeReports.value = []
	}
}

const loadReadme = async () => {
	readmeLoading.value = true
	readmeError.value = false
	readmeContent.value = ''

	try {
		// Fetch README from API
		const readmeData = await apiStore.fetchPackageReadme(pkg.value.id)

		if (readmeData && readmeData.content) {
			readmeContent.value = readmeData.content
		} else {
			readmeError.value = true
		}
	} catch (error) {
		console.error('Failed to load README:', error)
		readmeError.value = true
	} finally {
		readmeLoading.value = false
	}
}

const toggleBookmark = async () => {
	if (!userStore.isLoggedIn) {
		router.push('/login')
		return
	}

	bookmarkLoading.value = true
	try {
		if (pkg.value.is_bookmarked) {
			await apiStore.unbookmarkPackage(pkg.value.id)
			pkg.value.is_bookmarked = false
			pkg.value.bookmark_count--
		} else {
			await apiStore.bookmarkPackage(pkg.value.id)
			pkg.value.is_bookmarked = true
			pkg.value.bookmark_count++
		}

		Notify.create({
			type: 'positive',
			message: pkg.value.is_bookmarked ? 'Package bookmarked' : 'Bookmark removed',
			position: 'top',
		})
	} catch (_error) {
		// Error already handled by apiStore
	} finally {
		bookmarkLoading.value = false
	}
}

const copyCloneCommand = async () => {
	try {
		await copyToClipboard(`git clone ${pkg.value.repository_url}`)
		Notify.create({
			type: 'positive',
			message: 'Copied to clipboard',
			position: 'top',
			timeout: 2000,
		})
	} catch (error) {
		console.error('Failed to copy:', error)
	}
}

const getRepoType = (url) => {
	if (url.includes('github.com')) return 'GitHub'
	if (url.includes('gitlab.com')) return 'GitLab'
	if (url.includes('bitbucket.org')) return 'Bitbucket'
	return 'Git'
}

const getReportSeverityColor = (reason) => {
	const severities = {
		'Malicious code': 'red-10',
		'Copyright violation': 'deep-orange-10',
		'Inappropriate content': 'orange-10',
		'Broken/non-functional': 'amber-10',
		Spam: 'grey-10',
		Other: 'blue-grey-10',
	}
	return severities[reason] || 'grey-10'
}

const getReputationIcon = (rank) => {
	const icons = {
		banned: 'block',
		probation: 'warning',
		neutral: 'person',
		trusted: 'verified_user',
		contributor: 'star',
		veteran: 'military_tech',
	}
	return icons[rank] || 'person'
}

const getReputationColor = (rank) => {
	const colors = {
		banned: 'negative',
		probation: 'warning',
		neutral: 'grey',
		trusted: 'primary',
		contributor: 'positive',
		veteran: 'amber',
	}
	return colors[rank] || 'grey'
}

const getTagColor = (tag) => {
	if (tag.net_score >= 5) return 'primary'
	if (tag.net_score > 0) return 'grey-7'
	return 'grey-5'
}

const voteOnTag = async (tag, direction) => {
	if (!userStore.isLoggedIn) {
		router.push('/login')
		return
	}

	// Toggle vote if clicking same direction
	const newVote = tag.user_vote === direction ? 0 : direction
	const voteChange = newVote - tag.user_vote

	// Optimistic update
	tag.voting = true
	tag.user_vote = newVote
	tag.net_score = (tag.net_score || 0) + voteChange

	try {
		const result = await apiStore.voteOnTag(pkg.value.id, tag.id, newVote)

		// Update with actual server values
		if (result) {
			tag.net_score = result.net_score

			// If tag was removed due to low score, remove it from the UI
			if (result.removed) {
				const index = pkg.value.tags.findIndex((t) => t.id === tag.id)
				if (index !== -1) {
					pkg.value.tags.splice(index, 1)
				}
				Notify.create({
					type: 'info',
					message: 'Tag removed due to negative votes',
					position: 'top',
					timeout: 2000,
				})
				return
			}
		}

		Notify.create({
			type: 'positive',
			message: newVote === 0 ? 'Vote removed' : `Voted ${direction > 0 ? '+' : '-'}1 on tag`,
			position: 'top',
			timeout: 2000,
		})
	} catch (error) {
		// Revert on error
		tag.user_vote = direction === 0 ? 0 : -direction
		tag.net_score -= voteChange
	} finally {
		tag.voting = false
	}
}

const submitFlag = async () => {
	if (!flagReason.value) {
		Notify.create({
			type: 'negative',
			message: 'Please select a reason',
			position: 'top',
		})
		return
	}

	flagLoading.value = true
	try {
		await apiStore.flagPackage(pkg.value.id, flagReason.value, flagDetails.value)

		showFlagDialog.value = false

		// Reload flags to show the new one
		const flags = await apiStore.fetchPackageFlags(pkg.value.id)
		activeReports.value = flags.map((flag) => ({
			...flag,
			retracting: false,
		}))
	} catch (error) {
		console.error('Failed to submit flag:', error)
		Notify.create({
			type: 'negative',
			message: 'Failed to submit report',
			position: 'top',
		})
	} finally {
		flagLoading.value = false
	}
}

const retractReport = async (flagId) => {
	const report = activeReports.value.find((r) => r.id === flagId)
	if (!report) return

	report.retracting = true
	try {
		await apiStore.deleteFlag(flagId)

		// Remove from UI
		activeReports.value = activeReports.value.filter((r) => r.id !== flagId)
	} catch (error) {
		console.error('Failed to retract report:', error)
		Notify.create({
			type: 'negative',
			message: 'Failed to retract report',
			position: 'top',
		})
	} finally {
		report.retracting = false
	}
}

const filterTags = async (val, update) => {
	// For now, just allow any tag name
	update(() => {
		availableTags.value = []
	})
}

const createTag = (val, done) => {
	// Allow any tag name
	done(val)
}

const addTag = async () => {
	if (!newTagName.value || !newTagName.value.trim()) return

	addingTag.value = true
	try {
		const result = await apiStore.addPackageTag(pkg.value.id, newTagName.value)

		if (result) {
			// Create new tag object with initial vote
			const newTag = {
				id: result.tag_id,
				name: result.tag_name,
				net_score: result.vote_value,
				user_vote: 1,
				voting: false,
			}

			// Check if tag already exists before adding
			if (!pkg.value.tags) {
				pkg.value.tags = []
			}

			const existingTag = pkg.value.tags.find((t) => t.id === result.tag_id)
			if (!existingTag) {
				pkg.value.tags.push(newTag)
			}

			Notify.create({
				type: 'positive',
				message: 'Tag added successfully',
				position: 'top',
			})

			// Always close dialog and clear input on success
			showAddTagDialog.value = false
			newTagName.value = ''
		}
	} catch (error) {
		console.error('Failed to add tag:', error)
		// Still close the dialog on error to avoid stuck state
		showAddTagDialog.value = false
		newTagName.value = ''
	} finally {
		addingTag.value = false
	}
}

onMounted(() => {
	loadPackage()
})
</script>

<style lang="scss" scoped>
.clone-command {
	font-family: monospace;
}

.tag-with-voting {
	display: inline-flex;
	align-items: center;
	position: relative;

	&:hover .tag-voting-buttons {
		opacity: 1;
	}
}

.tag-voting-buttons {
	display: flex;
	margin-left: 4px;
	opacity: 0;
	transition: opacity 0.2s;

	.q-btn {
		padding: 2px;
	}
}
</style>
