<script lang="ts">
	import Tags from '$components/tags.svelte';
	import { timeAgo } from '$lib/utils';

	export let pkg: App.SearchResult;
	let packageName = pkg.name.toLowerCase().replace(/\s/g, '-');
</script>

<div class="search-result">
	<!-- <pre>{JSON.stringify(pkg, null, 2)}</pre> -->
	<h2 class="search-result--header">
		<a href={`/packages/${pkg.owner}/${packageName}`}>
			<span class="color-primary">{pkg.name}</span>
			<span class="result-path">( {pkg.owner}/{pkg.slug} )</span>
		</a>
	</h2>
	<span class="block">{pkg.description}</span>
	<!-- <pre>{JSON.stringify(pkg, null, 2)}</pre> -->
	<span class="block">
		Imports: <span class="color-primary">{pkg.dependency_count ?? '0'}</span> | v{pkg.version} @ {timeAgo(
			new Date(pkg.last_updated)
		)} |
		<span class="color-primary">{pkg.license}</span>
		|
		<Tags tags={pkg.keywords} />
	</span>
	<!-- <p>copy to clipboard</p> -->
</div>

<style>
	h2 {
		padding: 0;
		margin: 0;
	}
	.search-result {
		display: flex;
		flex-direction: column;
		gap: 0.1rem;
		padding: 0 0 2rem;
	}
	.search-result--header {
		align-items: center;
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		font-size: 1.3rem;
	}

	.result-path {
		color: var(--color-text);
	}
	.block {
		display: block;
	}
</style>
