<script lang="ts">
	import DetailsHeader from '$components/package/detailsHeader.svelte';
	import TabsHeader from '$components/package/tabsHeader.svelte';
	import Readme from '$components/package/readme.svelte';
	import Signatures from '$components/package/signatures.svelte';
	import Dependencies from '$components/package/dependencies.svelte';
	import { timeAgo } from '$lib/utils';
	//
	export let details: App.PackageDetails;
	export let readme;
	export let versionIndex;
	export let flat;
	export let licenses;
	let selectedTab: number = 0;
	$: details, (selectedTab = 0);
	$: readme = readme ?? 'No Readme Provided';

	function handleTabSelect(event: CustomEvent) {
		selectedTab = event.detail;
	}
	$: depCount = details.versions[versionIndex].dependency_count;
</script>

<svelte:head>
	<title>{details.repo_name} | OPM</title>
</svelte:head>

<main>
	<DetailsHeader {details} {versionIndex} on:goto_deps={() => (selectedTab = 2)} />
	<TabsHeader
		tabs={[
			{ name: 'Readme' },
			{ name: 'Signatures', disabled: true },
			{ name: `Dependencies (${depCount})`, disabled: depCount == 0 },
			{ name: `Versions (${details.versions.length})` }
		]}
		{selectedTab}
		on:select={handleTabSelect}
	/>
	{#if selectedTab === 0}
		<Readme {readme} />
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
					<th>Commit Hash</th>
					<!-- <th>Direct Security Issue?</th>
					<th>Indirect Security Issue?</th> -->
				</tr>
			</thead>
			<tbody>
				{#each details.versions as info}
					<tr>
						<td>
							<a href="/{details.host_name}/{details.owner_name}/{details.repo_name}?version={info.version}">
								{info.version}
							</a>
						</td>
						<td>{timeAgo(info.createdat)}</td>
						<td>{info.size_kb}</td>
						<td>{info.dependency_count}</td>
						<td>{info.compiler}</td>
						<td>{info.license}</td>
						<td class="mono">{info.commit_hash}</td>
						<!-- <td>{info.insecure}</td>
						<td>{info.has_insecure_dependency}</td> -->
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</main>

<style>
	table {
		max-width: 90vw;
		width: 100%;
		text-align: center;
	}
	.mono {
		font-family: 'Courier New', Courier, monospace;
	}
</style>
