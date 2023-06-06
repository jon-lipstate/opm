<script lang="ts">
	import { page } from '$app/stores';
	import SearchResult from '$components/searchResult.svelte';
	import Pagination from '$components/pagination.svelte';
	import { goto } from '$app/navigation';

	$: offset = Number($page.url.searchParams.get('offset')) ?? 0;
	$: results = $page.data.results.values as App.SearchResult[];
	$: count = $page.data.results.count;
	$: currentPage = Math.floor(offset / 100) + 1; // 100 results per page
	$: totalPages = Math.ceil(count / 100);
	$: {
		// Reactive Page Update:
		if (typeof window !== 'undefined') {
			currentPage;
			const newOffset = (currentPage - 1) * 100;
			if (newOffset != 0) {
				goto(`/browse?offset=${newOffset}`);
			}
		}
	}
</script>

<svelte:head>
	<title>Browse | Page {currentPage}</title>
	<meta name="description" content="Browse Page" />
</svelte:head>

<section>
	<h2>Total Packages: {count}</h2>
	{#each results as pkg (pkg.host_name + pkg.owner_name + pkg.repo_name)}
		<SearchResult {pkg} />
	{/each}
	{#if totalPages > 1}
		<Pagination bind:currentPage {totalPages} />
	{/if}
</section>
