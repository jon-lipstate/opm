<script lang="ts">
	import { page } from '$app/stores';
	import Package from '$components/package/_package.svelte';
	import type { LicenseSummary, FlatDependencies } from 'src/routes/api/details/dependencies/+server';
	// TODO: Freestanding Badge? Non-portable?
	//@ts-ignore
	export let data = $page.data;
	$: details = data.details;
	$: readme = data.readme.readme;
	$: versionIndex = data.versionIndex;
	$: flat = data.flat as FlatDependencies[];
	$: licenses = data.licenses as LicenseSummary[];
	$: error = data.error;
</script>

<svelte:head>
	<title>OPM: {details.repo_name ?? 'Invalid Package'}</title>
</svelte:head>

{#if error.length != 0}
	<pre>{error}</pre>
{/if}
<!-- <pre>{JSON.stringify(details, null, 2)}</pre> -->
{#if details}
	<Package {details} {readme} {versionIndex} {flat} {licenses} />
{:else}
	<h3>{error}</h3>
{/if}

<style>
	pre {
		font-size: 1.5rem;
		color: red;
	}
</style>
