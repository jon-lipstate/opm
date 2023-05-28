<script>
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	//
	export let currentPage = 1;
	export let totalPages = 1;
	let neighbors = 2;

	const navigate = (page) => {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
			dispatch('navigate', { page });
		}
	};
	$: currentPage, navigate(currentPage);
</script>

<nav>
	<button disabled={currentPage === 1} on:click={() => currentPage--}>&lt;</button>

	{#if currentPage - neighbors > 1}
		<button on:click={() => (currentPage = 1)}>1</button>
		<span>...</span>
	{/if}

	{#each Array.from({ length: neighbors * 2 + 1 }, (_, i) => currentPage - neighbors + i) as page}
		{#if page > 0 && page <= totalPages}
			<button class={page === currentPage ? 'active' : ''} on:click={() => (currentPage = page)}>{page}</button>
		{/if}
	{/each}

	{#if currentPage + neighbors < totalPages}
		<span>...</span>
		<button on:click={() => (currentPage = totalPages)}>{totalPages}</button>
	{/if}

	<button disabled={currentPage === totalPages} on:click={() => currentPage++}>&gt;</button>
</nav>

<style>
	nav {
		text-align: center;
	}
	.active {
		font-weight: bold;
		color: var(--color-theme-1);
	}
	button {
		appearance: none;
		background: none;
		border: none;
		color: inherit;
		cursor: pointer;
		padding: 0 0.15rem;
		font: inherit;
	}

	button:disabled {
		color: gray;
		cursor: not-allowed;
	}

	button:hover {
		text-decoration: underline;
	}
</style>
