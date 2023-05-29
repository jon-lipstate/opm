<script lang="ts">
	import { page } from '$app/stores';
	//
	import Pagination from '$components/pagination.svelte';
	import SearchResult from '$components/searchResult.svelte';
	import SearchBar from '$components/searchBar.svelte';

	$: query = $page.url.searchParams.get('query');
	$: queryMsg = query ?? 'No query provided';
	$: results = $page.data.results as App.SearchResult[];

	$: offset = Number($page.url.searchParams.get('offset')) ?? 0;
	$: count = results?.length ?? 0;
	$: currentPage = Math.floor(offset / 100) + 1; // 100 results per page
	$: totalPages = Math.ceil(count / 100);
</script>

<svelte:head>
	<title>Search Results for '{queryMsg}'</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<SearchBar />
	<h1>Results for '{queryMsg}' ({count} results)</h1>
	{#each results as pkg}
		<SearchResult {pkg} />
	{/each}
	{#if totalPages > 1}
		<Pagination bind:currentPage {totalPages} />
	{/if}
</section>

<style>
</style>
