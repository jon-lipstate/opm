<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import axios from 'axios';
	import Package from '$components/package/_package.svelte';
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
</script>

<svelte:head>
	<title>OPM: {details.name}</title>
</svelte:head>

<Package {details} />

<style>
</style>
