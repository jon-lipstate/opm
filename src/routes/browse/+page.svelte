<script lang="ts">
	import axios from 'axios';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import SearchResult from '$components/searchResult.svelte';
	import Pagination from '$components/pagination.svelte';
	import { goto } from '$app/navigation';

	$: offset = Number($page.url.searchParams.get('offset')) ?? 0;
	$: results = [] as App.SearchResult[];
	$: count = results.length;
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
	onMount(async () => {
		const response = await axios.post(`api/browse`, { offset });
		results = response.data.values;
	});
</script>

<svelte:head>
	<title>Browse | Page {currentPage + 1}</title>
	<meta name="description" content="Browse Page" />
</svelte:head>

<section>
	<span>Total Packages: {count}</span>
	{#each results as pkg (pkg.name)}
		<SearchResult {pkg} />
	{/each}
	{#if totalPages > 1}
		<Pagination bind:currentPage {totalPages} />
	{/if}
</section>
