<script>
	import { goto } from '$app/navigation';

	let query = '';

	async function handleSearch(e) {
		try {
			goto(`/search?query=${encodeURIComponent(query)}`);
		} catch (error) {
			console.error('handleSearch Error:', error);
		}
	}

	if (typeof window !== 'undefined') {
		window.addEventListener('keydown', function (e) {
			if (e.ctrlKey && e.key === '/') {
				e.preventDefault();
				//@ts-ignore
				document.getElementById('query').focus();
			}
		});
	}
</script>

<div class="search-container">
	<form action="api/search" on:submit|preventDefault={handleSearch}>
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
		<button type="submit" style="display: none;" />
		<!-- Hidden submit button to trigger form submission -->
		<div class="shortcut">
			<kbd class="">Ctrl</kbd> <kbd class="">/</kbd>
		</div>
	</form>
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<span class="advanced-menu" on:click={() => alert('not impl')}>Advanced</span>
</div>

<style>
	.advanced-menu {
		margin: 5px;
		color: var(--color-theme-4);
		cursor: pointer;
	}
	.search-container {
		position: relative;
		margin: 1rem 0;
		padding: 0 3px;
	}
	/* hide the x in the search: */
	input[type='search']::-webkit-search-cancel-button {
		display: none;
	}
	.search-input {
		padding: 0.5rem 0.5rem;
		border: 1px solid var(--back-translucent);
		/* https://font.gooova.com/fonts/13876/ */
		font-family: sans-serif;
		font-size: 1.5rem;
		appearance: none;
		width: 100%;
		height: 2.5rem;
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
