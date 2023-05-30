<script lang="ts">
	import { page } from '$app/stores';
	import DetailsHeader from '$components/package/detailsHeader.svelte';
	import TabsHeader from '$components/package/tabsHeader.svelte';
	import Readme from '$components/package/readme.svelte';
	import Signatures from '$components/package/signatures.svelte';
	import Dependencies from '$components/package/dependencies.svelte';
	import { timeAgo } from '$lib/utils';
	//
	export let details: App.PackageDetails;
	export let versionIndex;
	export let flat;
	export let licenses;
	let selectedTab: number = 0;
	$: details, (selectedTab = 0);

	let readmeData = 'no readme available';
	function handleTabSelect(event: CustomEvent) {
		selectedTab = event.detail;
	}
	$: depCount = details.versions[versionIndex].dependency_count;
</script>

<svelte:head>
	<title>{details.name} | OPM</title>
</svelte:head>

<main>
	<DetailsHeader {details} {versionIndex} on:goto_deps={() => (selectedTab = 2)} />
	<TabsHeader
		tabs={[
			{ name: 'Readme' },
			{ name: 'Signatures' },
			{ name: `Dependencies (${depCount})`, disabled: depCount == 0 },
			{ name: `Versions (${details.versions.length})` }
		]}
		{selectedTab}
		on:select={handleTabSelect}
	/>
	{#if selectedTab === 0}
		<!-- todo: sanitize -->
		<div>{@html readmeData}</div>
	{:else if selectedTab === 1}
		<Signatures />
	{:else if selectedTab === 2}
		<Dependencies {flat} {licenses} />
	{:else if selectedTab === 3}
		<table>
			<thead>
				<tr>
					<th>Version</th>
					<th>Date</th>
					<th>Size (kb)</th>
					<th>Imports</th>
					<th>Compiler</th>
					<th>License</th>
					<th>Direct Security Issue?</th>
					<th>Indirect Security Issue?</th>
				</tr>
			</thead>
			<tbody>
				{#each details.versions as info}
					<tr>
						<td><a href="/packages/{details.owner}/{details.slug}?version={info.version}">{info.version}</a></td>
						<td>{timeAgo(info.createdat)}</td>
						<td>{info.size_kb}</td>
						<td>{info.dependency_count}</td>
						<td>{info.compiler}</td>
						<td>{info.license}</td>
						<td>{info.insecure}</td>
						<td>{info.has_insecure_dependency}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</main>

<style>
	table {
		width: 100%;
		text-align: center;
	}
</style>
