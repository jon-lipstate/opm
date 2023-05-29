<script lang="ts">
	import { page } from '$app/stores';
	import Package from '$components/package/_package.svelte';
	// TODO: Freestanding Badge? Non-portable?
	let details = $page.data.details;
</script>

<svelte:head>
	<title>OPM: {details.name}</title>
</svelte:head>

<Package {details} />

<style>
</style>
