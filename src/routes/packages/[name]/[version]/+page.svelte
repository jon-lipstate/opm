<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import DetailsHeader from './detailsHeader.svelte';
	import TabsHeader from './tabsHeader.svelte';
	import Readme from './readme.svelte';
	import Signatures from './signatures.svelte';
	import Dependencies from './dependencies.svelte';
	import VersionRow from './versionRow.svelte';
	import axios from 'axios';
	// TODO: Freestanding Badge? Non-portable?
	let details = {
		name: '',
		version: '',
		description: '',
		tags: [],
		versions: [],
		funding: [],
		dependsOn: [],
		usedBy: [],
		requirements: { minCompilierVersion: '' },
		links: { repo: '', discord: '' },
		lastUpdated: '',
		license: '',
		size: '',
		kind: '',
		owners: [{ name: '', username: '' }],
		stats: { allTimeDownloads: 0 },
		readme: ''
	};

	let versions = [];

	onMount(async () => {
		let name = $page.params.name;
		let version = $page.params.version;
		// Fetch data from the API
		let response = await axios.post(`/api/details`, { name, version });
		// let response = await fetch(
		// 	`/api/details?package=${encodeURIComponent(name)}&version=${encodeURIComponent(version)}`
		// );
		if (response.status == 200) {
			details = response.data;
			versions = details.versions;
		} else {
			console.error(response.status, response.statusText);
		}
	});

	let selectedTab: number = 0;
	function handleTabSelect(event: CustomEvent) {
		selectedTab = event.detail;
	}
</script>

<svelte:head>
	<title>OPM: {details.name}</title>
</svelte:head>

<main>
	<DetailsHeader {details} />
	<TabsHeader tabs={['Readme', 'Signatures', 'Dependencies', 'Versions']} on:select={handleTabSelect} />

	<!-- <p>Currently selected tab: {selectedTab}</p> -->
	{#if selectedTab === 0}
		<Readme readme={details.readme} />
	{:else if selectedTab === 1}
		<Signatures />
	{:else if selectedTab === 2}
		<Dependencies />
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
				{#each versions as info}
					<VersionRow {info} />
				{/each}
			</tbody>
		</table>
	{/if}
</main>

<style>
</style>
