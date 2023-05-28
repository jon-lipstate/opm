<script lang="ts">
	import axios from 'axios';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	//
	import Pagination from '$components/pagination.svelte';
	import SearchResult from '$components/searchResult.svelte';
	import SearchBar from '$components/searchBar.svelte';

	$: query = $page.url.searchParams.get('query');
	$: queryMsg = query ?? 'No query provided';
	$: results = [] as App.SearchResult[];

	$: offset = Number($page.url.searchParams.get('offset')) ?? 0;
	$: count = results.length;
	$: currentPage = Math.floor(offset / 100) + 1; // 100 results per page
	$: totalPages = Math.ceil(count / 100);

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
	{#if totalPages > 1}
		<Pagination bind:currentPage {totalPages} />
	{/if}
</section>

<style>
</style>
