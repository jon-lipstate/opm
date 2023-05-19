<script lang="ts">
	import { page } from '$app/stores';
	import SearchResult from './searchResult.svelte';

	let query = $page.url.searchParams.get('query') || 'No query provided';

	const results = [
		{
			name: "Yuki's ECS",
			version: '1.2.3',
			updated: '1/1/1901',
			downloads: 100000,
			tags: ['ecs', 'engine'],
			kind: 'Curated Library'
		},
		{
			name: 'pico editor',
			version: '0.0.1',
			updated: '2/2/2902',
			downloads: 10,
			tags: ['tui', 'text-editor'],
			kind: 'demo'
		}
	];
</script>

<svelte:head>
	<title>Search Results for '{query}'</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<h1>Results for '{query}'</h1>
	{#each results as pkg}
		<SearchResult {pkg} />
	{/each}
</section>

<style>
</style>
