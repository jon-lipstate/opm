<script lang="ts">
	import Tags from '$components/tags.svelte';
	import leaf from '$lib/icons/leaf.svg';
	import yellowSand from '$lib/icons/yellow-sand.svg';
	// import redSand from '$lib/icons/red-sand.svg';
	import { isStale, timeAgo } from '$lib/utils';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	//
	export let details: App.PackageDetails;
	export let versionIndex = 0;

	$: version = details.versions[versionIndex];
</script>

<header class="details-header">
	<div class="row">
		<div>
			<h1 style="display:inline">
				<a href={details.url}>{details.repo_name}</a>
			</h1>
			<span>{details.host_name}/{details.owner_name}/{details.repo_name}</span>
		</div>
		{#if !version.insecure && version.has_insecure_dependency}
			<span class="insecure-warning">VULNERABLE DEPENDENCIES</span>
		{:else if version.insecure}
			<span class="insecure-warning">REPORTED VULNERABILITIES</span>
		{/if}
		{#if details.state == 'archived'}
			<span class="archived">ARCHIVED</span>
		{/if}
		<div class="features">
			<span>v{version.version}</span>
			|
			<span class="license">{version.license}</span>
			|
			<span>
				Updated: {timeAgo(version.createdat)}
				{#if isStale(version.createdat)}
					<img src={yellowSand} alt="stale" />
				{:else}
					<img src={leaf} alt="fresh" />
				{/if}
			</span>
			|
			<span>{version.size_kb} kb</span>

			<!-- <span>Used By: <a href="#">{details.usedBy?.length}</a></span> -->
		</div>
	</div>

	<div class="row">
		<!-- <a class="repo-link" href={details.links.url}> <img src={gh} alt="github logo" class="github-logo" />Repository</a> -->
		<Tags tags={details.keywords} />
		<div>
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			{#if version.dependency_count == 0}
				<span style="color:rgba(50,200,75,0.8)"> Dependencies: None </span>
			{:else}
				<span class="link" on:click={() => dispatch('goto_deps')}>
					Dependencies: {version.dependency_count}
				</span>
			{/if}
			|
			<span>Compiler: {version.compiler}</span>
		</div>
	</div>
</header>

<style>
	.archived {
		color: orange;
		font-weight: 900;
	}
	.insecure-warning {
		color: red;
		font-weight: 900;
	}
	.link {
		color: var(--color-theme-1);
	}
	.link:hover {
		text-decoration: underline;
		cursor: pointer;
	}
	.license {
		color: var(--c-odin-orange-lighten-3);
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
