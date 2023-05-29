<script lang="ts">
	import { page } from '$app/stores';
	import Package from '$components/package/_package.svelte';
	// TODO: Freestanding Badge? Non-portable?
	let details = $page.data.details;
	let error = $page.data.error;
</script>

<svelte:head>
	<title>OPM: {details?.name ?? 'Invalid Package'}</title>
</svelte:head>

{#if details}
	<!-- <pre>details {JSON.stringify(details, null, 2)}</pre> -->
	<Package {details} />
{:else}
	<h3>{error}</h3>
{/if}

<style>
</style>
