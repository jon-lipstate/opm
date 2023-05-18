<script>
</script>

<svelte:head>
	<title>Package Search</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<section>
	<!-- Search Bar-->
	<form class="search-container" action="/search">
		<input
			type="search"
			name="query"
			class="search-input"
			placeholder="Search"
			spellcheck="false"
			autocomplete="off"
		/>
		<div class="shortcut">
			<kbd class="">Ctrl</kbd> <kbd class="">K</kbd>
		</div>
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
	.search-input {
		padding: 0.2em 0.5em;
		border: 1px solid var(--back-translucent);
		font-family: inherit;
		font-size: 2rem;
		appearance: none;
		width: 100%;
		height: 1.3em;
		background-color: #3883d235;
		border-radius: var(--border-radius);
		vertical-align: middle;
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
