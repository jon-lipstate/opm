<script lang="ts">
	import axios from 'axios';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { searchResults } from '$stores/search';

	import SearchResult from '$components/searchResult.svelte';
	import SearchBar from '$components/searchBar.svelte';

	let results = [];
	searchResults.subscribe((value) => {
		results = value;
	});

	$: query = $page.url.searchParams.get('query');
	$: queryMsg = query ?? 'No query provided';

	onMount(async () => {
		//TODO: persist the results instead to not hit the server with another query?
		if ((!results || results.length == 0) && !!query) {
			const response = await axios.post(`api/search`, { query });
			results = response.data;
		}
	});
</script>

<svelte:head>
	<title>Search Results for '{queryMsg}'</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<SearchBar />
	<h1>Results for '{queryMsg}' ({results?.length ?? 0} results)</h1>
	{#each results as pkg}
		<SearchResult {pkg} />
	{/each}
</section>

<style>
</style>
