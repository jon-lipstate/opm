<script lang="ts">
	import Tags from '$components/tags.svelte';
	import gh from '$lib/images/github.svg';
	import leaf from '$lib/icons/leaf.svg';
	import yellowSand from '$lib/icons/yellow-sand.svg';
	// import redSand from '$lib/icons/red-sand.svg';
	import { timeAgo } from '$lib/utils';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	export let details: App.PackageDetails;
</script>

<header class="details-header">
	<div class="row">
		<h1>
			<a href={details.links.url}>{details.name}</a>
		</h1>
		<div class="features">
			<span>v{details.version}</span>
			|
			<span class="license">{details.license}</span>
			|
			<span>
				Updated: {timeAgo(details.lastUpdated)}
				{#if true}
					<img src={leaf} alt="fresh" />
				{:else}
					<img src={yellowSand} alt="stale" />
					<!-- <img src="{redSand}" alt="inactive"> -->
				{/if}
			</span>
			|
			<span>{details.size} kb</span>

			<!-- <span>Used By: <a href="#">{details.usedBy?.length}</a></span> -->
		</div>
	</div>

	<div class="row">
		<!-- <a class="repo-link" href={details.links.url}> <img src={gh} alt="github logo" class="github-logo" />Repository</a> -->
		<Tags tags={details.tags} />
		<div>
			<span>
				Depends On: <a href="#/" on:click={() => dispatch('clickDeps')}>
					{Object.keys(details.dependsOn)?.length ?? 0}
				</a>
			</span>
			|
			<span>Compiler: {details.requirements.compiler}</span>
		</div>
	</div>
</header>

<style>
	.license {
		color: greenyellow;
	}
	.features {
		margin-left: 1rem;
	}
	.details-header {
		color: var(--color-theme-4);
		background-color: rgba(194, 226, 229, 0.172);
		padding: 0.05rem 1rem 0.2rem 1rem;
		position: relative;
		overflow: hidden;
	}
	.row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	h1 {
		margin: 0;
	}
</style>
