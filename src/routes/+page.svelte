<script>
	import axios from 'axios';
	import { goto } from '$app/navigation';
    import { searchResults } from '$stores/search';

	let query='';

	async function handleSearch(e) {
        try {
            const response = await axios.post(`${import.meta.env.VITE_API_HOST}/search`, { query });
			$searchResults = response.data;
			goto(`/search?query=${encodeURIComponent(query)}`);
        } catch (error) {
            console.error('Error:', error);
        }
    }

	if (typeof window !== "undefined") {
	window.addEventListener('keydown', function(e) {
    if (e.ctrlKey && e.key === 'k') {
      e.preventDefault();
	  //@ts-ignore
      document.getElementById('query').focus();
    }
  });
}
</script>

<svelte:head>
	<title>Package Search</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<!-- Search Bar--> 
	<form class="search-container" action="/search" on:submit|preventDefault={handleSearch}>
		<input
			bind:value={query}
			type="search"
			id="query"
			name="query"
			class="search-input"
			placeholder="Search"
			spellcheck="false"
			autocomplete="off"
		/>
		<button type="submit" style="display: none;"></button> <!-- Hidden submit button to trigger form submission -->
		<div class="shortcut">
			<kbd class="">Ctrl</kbd> <kbd class="">K</kbd>
		</div>
		<!-- <div>Package Search</div> -->
	</form>
</section>

<style>
	section {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 70%;
		min-width: 30rem;
		max-width: 50%;
	}
	.search-container {
		position: relative;
		margin: auto;
		padding: 0 3px;
	}
	/* hide the x in the search: */
	input[type='search']::-webkit-search-cancel-button {
		display: none;
	}
	.search-input {
		padding: 0.2em 0.5em;
		border: 1px solid var(--back-translucent);
		/* https://font.gooova.com/fonts/13876/ */
		font-family: sans-serif;
		font-size: 2rem;
		appearance: none;
		width: 100%;
		height: 1.3em;
		background-color: var(--color-bg-0);
		color: var(--color-text);
		border-radius: var(--border-radius);
		vertical-align: middle;
	}
	::placeholder { 
		color: lightgrey;
		opacity: 1; 
	}

	.shortcut {
		position: absolute;
		top: calc(0.6rem);
		right: 1rem;
		width: 100%;
		text-align: right;
		pointer-events: none;
		font-size: 1.2rem;
		text-transform: uppercase;
	}
</style>
