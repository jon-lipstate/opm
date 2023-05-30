<script lang="ts">
	import { page } from '$app/stores';
	import Package from '$components/package/_package.svelte';
	import type { LicenseSummary, FlatDependencies } from 'src/routes/api/details/dependencies/+server';
	// TODO: Freestanding Badge? Non-portable?
	//@ts-ignore
	export let data = $page.data;
	$: details = data.details;
	$: versionIndex = data.versionIndex;
	$: flat = data.flat as FlatDependencies[];
	$: licenses = data.licenses as LicenseSummary[];
	$: error = data.error;
</script>

<svelte:head>
	<title>OPM: {details?.name ?? 'Invalid Package'}</title>
</svelte:head>

{#if details}
	<Package {details} {versionIndex} {flat} {licenses} />
{:else}
	<h3>{error}</h3>
{/if}

<style>
</style>
