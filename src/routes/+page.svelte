<script>
	import SearchBar from '$components/searchBar.svelte';
</script>

<svelte:head>
	<title>Package Search</title>
	<meta name="description" content="Search Page" />
</svelte:head>

<class class="landing-search">
	<SearchBar />
</class>

<style>
	.landing-search {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 70vw;
		min-width: 20rem;
	}
</style>
