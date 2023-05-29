<script lang="ts">
	import { page } from '$app/stores';
	import DetailsHeader from '$components/package/detailsHeader.svelte';
	import TabsHeader from '$components/package/tabsHeader.svelte';
	import Readme from '$components/package/readme.svelte';
	import Signatures from '$components/package/signatures.svelte';
	import Dependencies from '$components/package/dependencies.svelte';
	import VersionRow from '$components/package/versionRow.svelte';
	//
	export let details: App.PackageDetails;
	// TODO: GET FROM QUERY STRING IF AVAIL
	$: versionIndex = details.versions.reduce((max, celm, cidx) => {
		return celm.id > details.versions[max].id ? cidx : max;
	}, 0);
	let readmeData = 'no readme available';
	let selectedTab: number = 0;
	function handleTabSelect(event: CustomEvent) {
		selectedTab = event.detail;
	}
</script>

<svelte:head>
	<title>{details.name} | OPM</title>
</svelte:head>

<main>
	<DetailsHeader {details} {versionIndex} on:goto_deps={() => (selectedTab = 2)} />
	<TabsHeader tabs={['Readme', 'Signatures', 'Dependencies', 'Versions']} {selectedTab} on:select={handleTabSelect} />
	<div style="color:yellow">TODO: NAV FROM DEPS LINK DOESNT UPDATE THE PAGE</div>
	<div style="color:yellow">TODO: disable deps tab when none, same for sigs since not doing yet</div>
	<div style="color:yellow">TODO: Versions Section</div>
	{#if selectedTab === 0}
		<!-- todo: sanitize -->
		<div>{@html readmeData}</div>
	{:else if selectedTab === 1}
		<Signatures />
	{:else if selectedTab === 2}
		<Dependencies versionId={details.versions[versionIndex].id} />
	{:else if selectedTab === 3}
		<table>
			<thead>
				<tr>
					<th>Version</th>
					<th>Date</th>
					<th>Changes</th>
				</tr>
			</thead>
			<tbody>
				{#each details.versions as info}
					<!-- <VersionRow {info} /> -->
					<p>n</p>
				{/each}
			</tbody>
		</table>
	{/if}
</main>

<style>
</style>
