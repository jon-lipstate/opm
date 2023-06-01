<script lang="ts">
	// TODO: make the gist and repo seperate pages/components?
	import axios from 'axios';
	import Slider from '$components/slider.svelte';
	import Select from '$components/select.svelte';
	import Package from '$components/package/_package.svelte';
	import { onMount } from 'svelte';
	import { intoPackageDetails } from '$lib/utils';
	import { page } from '$app/stores';
	import json5 from 'json5';
	//
	let session = $page.data.session;
	let selectedIndex;
	let repositories = [];
	let showForks = false;
	let missingPackageFile = false;
	let pkgFile: App.ModPkg | null = null;
	//@ts-ignore
	let details: App.PackageDetails = null;

	//@ts-ignore
	$: visibleRepos = repositories.filter((x) => showForks || !x.fork);
	// onMount(async () => {
	// 	try {
	// 		const repoRes = await axios.get(`/api/github/getPublicRepos`);
	// 		//@ts-ignore
	// 		repositories = repoRes.data;
	// 	} catch (e) {
	// 		console.error('api call error', e);
	// 	}
	// });

	async function handleSelectLib() {
		missingPackageFile = false;
		selLib = visibleRepos[selectedIndex];
		//@ts-ignore
		const name = visibleRepos[selectedIndex].name;
		try {
			const contentsResponse = await axios.post(`/api/github/getPkgFile`, { name });
			pkgFile = json5.parse(contentsResponse.data);
			//@ts-ignore
			const readmeResponse = await axios.get(`https://api.github.com/repos/${session.user.name}/${name}/readme`);
			details = intoPackageDetails(pkgFile!, selLib, readmeResponse.data.download_url);
		} catch (e) {
			console.error(e);
			//@ts-ignore
			if (e.response.status == 404) {
				missingPackageFile = true;
			} else {
				console.error('api call error', e);
			}
		}
	}
	let selLib: any;
	async function submitPackageDetails(e) {
		e.preventDefault();
		try {
			const response = await axios.post('/api/packages/new', details);
			console.log(response.data);
		} catch (error) {
			console.error('Error submitting package details', error, details);
		}
	}
</script>

<form on:submit={submitPackageDetails}>
	<div class="my-1">
		<Select
			label={'Repo'}
			bind:value={selectedIndex}
			prefix={'Select Repo'}
			options={visibleRepos}
			renderAs={(x) => x.name}
			on:change={handleSelectLib}
		/>
		<Slider bind:checked={showForks} />
		<span>{showForks ? 'Forks Shown' : 'Forks Hidden'}</span>
	</div>

	<button type="submit" disabled={selectedIndex === '' || missingPackageFile}>Submit</button>
</form>
{#if missingPackageFile}
	<div style="color:red">Missing mod.pkg file</div>
{/if}
<h3>Submission Preview</h3>
{#if details && !missingPackageFile}
	<Package {details} />
	<!-- <pre>{JSON.stringify(details, null, 2)}</pre> -->
{/if}
